package kafka

import (
	"log"

	"github.com/IBM/sarama"

	"github.com/bitmaskit/notifications/model"
)

type Kafka interface {
	BrokerAddr() string
	ProduceNotification(message model.NotificationRequest) error
	ProduceToTopic(msg string, topic string) error
}

type kafka struct {
	brokerAddress     string
	notificationTopic string
}

func New(brokerAddress string, notificationTopic string) Kafka {
	return &kafka{
		brokerAddress:     brokerAddress,
		notificationTopic: notificationTopic,
	}
}

func (k *kafka) BrokerAddr() string {
	return k.brokerAddress
}

func (k *kafka) ProduceNotification(msg model.NotificationRequest) error {
	jsonStr, err := msg.ToJSONString()
	if err != nil {
		log.Println("Error marshalling message: ", err)
		return err
	}
	producer, err := sarama.NewSyncProducer([]string{k.BrokerAddr()}, nil)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalf("Error closing producer: %v", err)
		}
	}()
	m := &sarama.ProducerMessage{
		Topic: k.notificationTopic,
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

func (k *kafka) ProduceToTopic(msg string, topic string) error {
	producer, err := sarama.NewSyncProducer([]string{k.BrokerAddr()}, nil)
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
		Value: sarama.StringEncoder(msg),
	}
	partition, offset, err := producer.SendMessage(m)
	if err != nil {
		log.Printf("Failed to produce message: %s", err)
		return err
	}
	log.Printf("Produced message to partition %d with offset %d\n", partition, offset)

	return nil
}
