package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"

	"github.com/bitmaskit/notifications/model"
	"github.com/bitmaskit/notifications/slack/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	log.Println("Starting slack service... Listening for notifications")
	consumeSlack(cfg)
}

func consumeSlack(cfg *config.SlackConfig) {
	// consumer code
	consumer, err := sarama.NewConsumer([]string{cfg.BrokerAddr}, nil)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalf("Error closing client: %v", err)
		}
	}()

	partitions, err := consumer.Partitions(cfg.SlackTopic)
	if err != nil {
		log.Fatalf("Error retrieving partitions: %v", err)
	}

	var wg sync.WaitGroup
	for _, partition := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(cfg.SlackTopic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Error consuming partition: %v", err)
		}
		defer func(consumer sarama.PartitionConsumer) {
			if err := consumer.Close(); err != nil {
				log.Fatalf("Error closing client: %v", err)
			}
		}(partitionConsumer)

		wg.Add(1)
		go func(partitionConsumer sarama.PartitionConsumer) {
			defer wg.Done()
			for {
				select {
				case msg := <-partitionConsumer.Messages():
					if err := sendToSlack(msg.Value, cfg); err != nil {
						log.Println("Error routing message:", err)
					}
				}
			}
		}(partitionConsumer)
	}

	// graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals
	wg.Wait()
}

func sendToSlack(value []byte, cfg *config.SlackConfig) error {
	sm := model.SlackMessage{
		Text: string(value),
	}
	msgBytes, err := json.Marshal(sm)
	if err != nil {
		return err
	}

	resp, err := http.Post(cfg.WebhookURL, "application/json", bytes.NewBuffer(msgBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message to Slack, status code: %d", resp.StatusCode)
	}

	log.Printf("Sending message to slack: %s\n", string(value))
	return nil
}
