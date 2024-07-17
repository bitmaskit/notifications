package main

import (
	"log"

	"github.com/bitmaskit/notifications/kafka"
	"github.com/bitmaskit/notifications/router/config"
	"github.com/bitmaskit/notifications/router/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	kafka := kafka.New(
		cfg.BrokerAddr,
		cfg.NotificationTopic,
	)
	log.Println("Starting router... Listening for notifications")
	router.ConsumeNotifications(kafka, cfg)
}
