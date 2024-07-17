package config

import (
	"errors"
	"os"
)

var (
	ErrBrokerAddrNotSet = errors.New("KAFKA_BROKER_ADDRESS is not set")
	ErrSlackTopicNotSet = errors.New("SLACK_KAFKA_TOPIC is not set")
	ErrWebhookUrlNotSet = errors.New("SLACK_WEBHOOK_URL is not set")
)

type SlackConfig struct {
	BrokerAddr string
	SlackTopic string
	WebhookURL string
}

func Load() (*SlackConfig, error) {
	brokerAddr := os.Getenv("KAFKA_BROKER_ADDRESS")
	if brokerAddr == "" {
		return nil, ErrBrokerAddrNotSet
	}
	slackTopic := os.Getenv("SLACK_KAFKA_TOPIC")
	if slackTopic == "" {
		return nil, ErrSlackTopicNotSet
	}
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	if webhookURL == "" {
		return nil, ErrWebhookUrlNotSet
	}

	return &SlackConfig{
		WebhookURL: webhookURL,
		BrokerAddr: brokerAddr,
		SlackTopic: slackTopic,
	}, nil
}
