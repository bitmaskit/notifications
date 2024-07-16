package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bitmaskit/notifications/channel"
	"github.com/bitmaskit/notifications/model"
	"github.com/bitmaskit/notifications/router/router"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

const (
	env   = ".env"
	topic = "notifications"
)

var brokerAddr string

func init() {
	if err := godotenv.Load(env); err != nil {
		log.Fatalf("Failed to load env: %v", err)
	}

	if brokerAddr = os.Getenv("KAFKA_BROKER_ADDRESS"); brokerAddr == "" {
		log.Fatalln("KAFKA_BROKER_ADDRESS is not set")
	}
}

func main() {
	// 1. consumer from notifications topic
	// 2. figure out which channels the message should go to
	// 3. send the message to the appropriate channel
	log.Println("Starting router... Listening for notifications")
	consumeNotifications()
}

func consumeNotifications() {
	// consumer code
	consumer, err := sarama.NewConsumer([]string{brokerAddr}, nil)
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
			fmt.Printf("Consumed message offset %d\n", msg.Offset)
			fmt.Printf("Message value: %s\n", string(msg.Value))
			message, channels, err := decodeMessage(msg.Value)
			if err != nil {
				log.Println("Error decoding message:", err)
			}
			if err := router.Route(message, channels); err != nil {
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

	log.Println("Decoded message:", nr.Message)
	log.Println("Decoded channels:", nr.Channels)

	return nr.Message, nr.Channels, nil
}
