package config

import (
	"errors"
	"os"
)

var (
	ErrBackendPortNotSet       = errors.New("BACKEND_PORT is not set")
	ErrKafkaBrokerAddrNotSet   = errors.New("KAFKA_BROKER_ADDRESS is not set")
	ErrNotificationTopicNotSet = errors.New("NOTIFICATIONS_KAFKA_TOPIC is not set")
)

type BackendConfig struct {
	Port              string
	BrokerAddr        string
	NotificationTopic string
}

func Load() (*BackendConfig, error) {
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		return nil, ErrBackendPortNotSet
	}
	brokerAddr := os.Getenv("KAFKA_BROKER_ADDRESS")
	if brokerAddr == "" {
		return nil, ErrKafkaBrokerAddrNotSet
	}
	notificationTopic := os.Getenv("NOTIFICATIONS_KAFKA_TOPIC")
	if notificationTopic == "" {
		return nil, ErrNotificationTopicNotSet
	}

	return &BackendConfig{
		Port:              port,
		BrokerAddr:        brokerAddr,
		NotificationTopic: notificationTopic,
	}, nil
}
