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

var (
	port       string
	brokerAddr string
)

func init() {
	if err := godotenv.Load(env); err != nil {
		log.Fatalf("Failed to load env: %v", err)
	}
	if port = os.Getenv("BACKEND_PORT"); port == "" {
		log.Fatalln("BACKEND_PORT is not set")
	}
	if brokerAddr = os.Getenv("KAFKA_BROKER_ADDRESS"); brokerAddr == "" {
		log.Fatalln("KAFKA_BROKER_ADDRESS is not set")
	}
}

func main() {
	api := api.API{
		Kafka: kafka.New(brokerAddr),
	}

	http.HandleFunc("/message", api.PostHandler)

	log.Println("Starting server on port: ", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalln("Error running the server", err)
	}
}
