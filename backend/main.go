package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/bitmaskit/notifications/backend/api"
	"github.com/bitmaskit/notifications/kafka"
)

const env = ".env"

var port string

func main() {
	if err := godotenv.Load(env); err != nil {
		log.Fatalf("Failed to load env: %v", err)
	}
	if port = os.Getenv("BACKEND_PORT"); port == "" {
		port = "8080"
	}

	api := api.API{Kafka: kafka.Kafka{}}

	http.HandleFunc("/message", api.PostHandler)

	log.Println("Starting server on port: ", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalln("Error running the server", err)
	}
}
