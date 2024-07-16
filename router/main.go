package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/bitmaskit/notifications/channel"
	"github.com/bitmaskit/notifications/kafka"
	"github.com/bitmaskit/notifications/model"
	"github.com/bitmaskit/notifications/router/router"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

const (
	env = ".env"
)

var (
	brokerAddr         string
	notificationsTopic string
)

func init() {
	if err := godotenv.Load(env); err != nil {
		log.Fatalf("Failed to load env: %v", err)
	}

	if brokerAddr = os.Getenv("KAFKA_BROKER_ADDRESS"); brokerAddr == "" {
		log.Fatalln("KAFKA_BROKER_ADDRESS is not set")
	}
	if notificationsTopic = os.Getenv("NOTIFICATIONS_TOPIC"); notificationsTopic == "" {
		log.Fatalln("NOTIFICATIONS_TOPIC is not set")
	}
}

func main() {
	log.Println("Starting router... Listening for notifications")
	kafka := kafka.New(brokerAddr)
	consumeNotifications(kafka, notificationsTopic)
}

func consumeNotifications(kafka kafka.Kafka, topic string) {
	// consumer code
	consumer, err := sarama.NewConsumer([]string{kafka.BrokerAddr()}, nil)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
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
			if err := router.Route(kafka, message, channels); err != nil {
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
