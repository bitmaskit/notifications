package config

type RouterConfig struct {
	BrokerAddr        string
	NotificationTopic string
	SmsTopic          string
	EmailTopic        string
	SlackTopic        string
}
