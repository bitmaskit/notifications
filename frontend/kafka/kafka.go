package kafka

import (
	"github.com/bitmaskit/notifications/channel"
	"log"
)

type Message struct {
	Message  string
	Channels []channel.Channel
}

type Kafka struct {
}

func (k Kafka) Produce(msg Message) error {
	// Produce message to kafka

	log.Println("Producing message: ", msg.Message)
	log.Println("Channels: ", msg.Channels)

	return nil
}
