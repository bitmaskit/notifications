package config

import (
	"errors"
	"os"
)

var (
	ErrBrokerAddrNotSet         = errors.New("KAFKA_BROKER_ADDRESS is not set")
	ErrSMSTopicNotSet           = errors.New("SMS_KAFKA_TOPIC is not set")
	ErrEmailTopicNotSet         = errors.New("EMAIL_KAFKA_TOPIC is not set")
	ErrSlackTopicNotSet         = errors.New("SLACK_KAFKA_TOPIC is not set")
	ErrNotificationsTopicNotSet = errors.New("NOTIFICATIONS_KAFKA_TOPIC is not set")
)

type RouterConfig struct {
	BrokerAddr        string
	NotificationTopic string
	SmsTopic          string
	EmailTopic        string
	SlackTopic        string
}

func Load() (*RouterConfig, error) {
	brokerAddr := os.Getenv("KAFKA_BROKER_ADDRESS")
	if brokerAddr == "" {
		return nil, ErrBrokerAddrNotSet
	}
	notificationsTopic := os.Getenv("NOTIFICATIONS_KAFKA_TOPIC")
	if notificationsTopic == "" {
		return nil, ErrNotificationsTopicNotSet
	}
	smsTopic := os.Getenv("SMS_KAFKA_TOPIC")
	if smsTopic == "" {
		return nil, ErrSMSTopicNotSet
	}
	emailTopic := os.Getenv("EMAIL_KAFKA_TOPIC")
	if emailTopic == "" {
		return nil, ErrEmailTopicNotSet
	}
	slackTopic := os.Getenv("SLACK_KAFKA_TOPIC")
	if slackTopic == "" {
		return nil, ErrSlackTopicNotSet
	}

	return &RouterConfig{
		BrokerAddr:        brokerAddr,
		NotificationTopic: notificationsTopic,
		SmsTopic:          smsTopic,
		EmailTopic:        emailTopic,
		SlackTopic:        slackTopic,
	}, nil
}
