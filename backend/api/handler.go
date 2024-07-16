package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bitmaskit/notifications/kafka"
	"github.com/bitmaskit/notifications/model"
)

const (
	SuccessfulMessage   = "Message successfully received and processed"
	UnsuccessfulMessage = "Error processing message"
)

type API struct {
	Kafka kafka.Kafka
}

func (a *API) PostHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse json
	var message model.NotificationRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&message); err != nil {
		log.Println("Error decoding request body: ", err)
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	response := model.ApiResponse{}
	w.Header().Set("Content-Type", "application/json")

	// send to kafka
	if err := a.Kafka.Produce(message); err != nil {
		log.Println("Error sending message to kafka: ", err)
		http.Error(w, "Error producing message", http.StatusInternalServerError)
		response.Message = UnsuccessfulMessage
		response.Success = false
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("Error encoding response: ", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
		return
	}

	response.Message = SuccessfulMessage
	response.Success = true

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding response: ", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
