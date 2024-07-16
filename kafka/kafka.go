package kafka

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/bitmaskit/notifications/model"
)

const (
	brokerAddress = "localhost:9092"
	topic         = "notifications"
)

type Kafka struct {
}

func (k Kafka) Produce(msg model.NotificationRequest) error {
	jsonStr, err := msg.ToJSONString()
	if err != nil {
		log.Println("Error marshalling message	: ", err)
		return err
	}
	producer, err := sarama.NewSyncProducer([]string{brokerAddress}, nil)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalf("Error closing producer: %v", err)
		}
	}()
	m := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(jsonStr),
	}
	partition, offset, err := producer.SendMessage(m)
	if err != nil {
		log.Printf("Failed to produce message: %s", err)
		return err
	}
	log.Printf("Produced message to partition %d with offset %d\n", partition, offset)

	return nil
}
