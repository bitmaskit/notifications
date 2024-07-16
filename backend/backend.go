package main

import (
	"log"
	"net/http"

	"github.com/bitmaskit/notifications/backend/api"
	"github.com/bitmaskit/notifications/backend/config"
	"github.com/bitmaskit/notifications/kafka"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalln("Failed to load config: ", err)
	}
	api := api.API{
		Kafka: kafka.New(
			cfg.BrokerAddr,
			cfg.NotificationTopic,
		),
	}

	http.HandleFunc("/message", api.PostHandler)

	log.Println("Starting server on port: ", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatalln("Error running the server", err)
	}
}
