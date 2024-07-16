package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const env = ".env"

var port string

func main() {
	if err := godotenv.Load(env); err != nil {
		log.Fatalf("Failed to load env: %v", err)
	}

	api := &api.API{}

	http.HandleFunc("GET /", api.IndexHandler)
	http.HandleFunc("POST /", api.PostHandler)

	if port = os.Getenv("FRONTEND_PORT"); port == "" {
		port = "8080"
	}
	log.Println("Starting server on port: ", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalln("Error running the server", err)
	}
}
