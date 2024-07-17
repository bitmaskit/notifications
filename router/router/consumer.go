package router

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/bitmaskit/notifications/kafka"
	"github.com/bitmaskit/notifications/router/config"
)

func ConsumeNotifications(kafka kafka.Kafka, cfg *config.RouterConfig) {
	// consumer code
	consumer, err := sarama.NewConsumer([]string{kafka.BrokerAddr()}, nil)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalf("Error closing client: %v", err)
		}
	}()

	partitions, err := consumer.Partitions(cfg.NotificationTopic)
	if err != nil {
		log.Fatalf("Error retrieving partitions: %v", err)
	}

	var wg sync.WaitGroup
	for _, partition := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(cfg.NotificationTopic, partition, sarama.OffsetNewest)
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
					if err := route(kafka, msg.Value, cfg); err != nil {
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
