package main

import (
	"log"
	"net/http"

	"github.com/bitmaskit/notifications/backend/api"
	"github.com/bitmaskit/notifications/kafka"
)

func main() {
	//api := backend.API{kafka.Kafka{}}
	api := api.API{Kafka: kafka.Kafka{}}
	http.HandleFunc("/message", api.PostHandler)

	log.Println("Starting server on :8001")
	if err := http.ListenAndServe(":8001", nil); err != nil {
		log.Fatalln("Error starting server: ", err)
	}
}
