package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bitmaskit/notifications/frontend/api"
	"github.com/bitmaskit/notifications/frontend/config"

	"github.com/joho/godotenv"
)

const env = ".env"

var (
	backendHost string
	port        string
)

func init() {
	if err := godotenv.Load(env); err != nil {
		log.Fatalf("Failed to load env: %v", err)
	}
	if port = os.Getenv("FRONTEND_PORT"); port == "" {
		log.Fatalln("FRONTEND_PORT is not set")
	}
	if backendHost = os.Getenv("BACKEND_HOST"); backendHost == "" {
		log.Fatalln("BACKEND_HOST is not set")
	}
}

func main() {

	api := api.New(config.New(backendHost))

	http.HandleFunc("GET /", api.IndexHandler)
	http.HandleFunc("POST /", api.PostHandler)

	log.Println("Starting server on port: ", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalln("Error running the server", err)
	}
}
