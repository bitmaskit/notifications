package main

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/bitmaskit/notifications/channel"
	"github.com/bitmaskit/notifications/kafka"
	"github.com/bitmaskit/notifications/model"
	"github.com/bitmaskit/notifications/router/config"
	"github.com/bitmaskit/notifications/router/router"

	"github.com/IBM/sarama"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	kafka := kafka.New(
		cfg.BrokerAddr,
		cfg.NotificationTopic,
	)
	log.Println("Starting router... Listening for notifications")
	consumeNotifications(kafka, cfg)
}

func consumeNotifications(kafka kafka.Kafka, cfg *config.RouterConfig) {
	// consumer code
	consumer, err := sarama.NewConsumer([]string{kafka.BrokerAddr()}, nil)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	partitionConsumer, err := consumer.ConsumePartition(cfg.NotificationTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error consuming partition: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalf("Error closing client: %v", err)
		}
	}()
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			message, channels, err := decodeMessage(msg.Value)
			if err != nil {
				log.Println("Error decoding message:", err)
			}
			if err := router.Route(kafka, message, channels, cfg); err != nil {
				log.Println("Error routing message:", err)
			}
		}
	}
}

func decodeMessage(msgValue []byte) (string, []channel.Channel, error) {
	nr := model.NotificationRequest{}
	reader := bytes.NewReader(msgValue)
	if err := json.NewDecoder(reader).Decode(&nr); err != nil {
		log.Println("Error decoding request body:", err)
		return "", nil, err
	}

	return nr.Message, nr.Channels, nil
}
