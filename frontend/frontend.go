package main

import (
	"github.com/bitmaskit/notifications/frontend/api"
	"github.com/bitmaskit/notifications/frontend/config"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalln("Failed to load config: ", err)
	}

	api := api.New(cfg)

	http.HandleFunc("GET /", api.IndexHandler)
	http.HandleFunc("POST /", api.PostHandler)

	log.Println("Starting server on port: ", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatalln("Error running the server", err)
	}
}
