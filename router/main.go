package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/bitmaskit/notifications/channel"
	"github.com/bitmaskit/notifications/kafka"
	"github.com/bitmaskit/notifications/model"
	"github.com/bitmaskit/notifications/router/config"
	"github.com/bitmaskit/notifications/router/router"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

const env = ".env"

var routerConfig *config.RouterConfig

func init() {
	if err := godotenv.Load(env); err != nil {
		log.Fatalf("Failed to load env: %v", err)
	}
	var brokerAddr string
	var notificationsTopic, smsTopic, emailTopic, slackTopic string

	if brokerAddr = os.Getenv("KAFKA_BROKER_ADDRESS"); brokerAddr == "" {
		log.Fatalln("KAFKA_BROKER_ADDRESS is not set")
	}
	if notificationsTopic = os.Getenv("NOTIFICATIONS_KAFKA_TOPIC"); notificationsTopic == "" {
		log.Fatalln("NOTIFICATIONS_KAFKA_TOPIC is not set")
	}
	if smsTopic = os.Getenv("SMS_KAFKA_TOPIC"); smsTopic == "" {
		log.Fatalln("SMS_KAFKA_TOPIC is not set")
	}
	if emailTopic = os.Getenv("EMAIL_KAFKA_TOPIC"); emailTopic == "" {
		log.Fatalln("EMAIL_KAFKA_TOPIC is not set")
	}
	if slackTopic = os.Getenv("SLACK_KAFKA_TOPIC"); slackTopic == "" {
		log.Fatalln("SLACK_KAFKA_TOPIC is not set")
	}

	routerConfig = &config.RouterConfig{
		BrokerAddr:        brokerAddr,
		NotificationTopic: notificationsTopic,
		SmsTopic:          smsTopic,
		EmailTopic:        emailTopic,
		SlackTopic:        slackTopic,
	}
}

func main() {
	log.Println("Starting router... Listening for notifications")
	kafka := kafka.New(
		routerConfig.BrokerAddr,
		routerConfig.NotificationTopic,
	)
	consumeNotifications(kafka, routerConfig)
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
