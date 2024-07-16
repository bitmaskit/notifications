package api

import (
	"html/template"
	"log"
	"net/http"

	"github.com/bitmaskit/notifications/channel"
	"github.com/bitmaskit/notifications/kafka"
	"github.com/bitmaskit/notifications/model"
)

const ErrInternalServerError = "Internal server error"

type API struct {
	Kafka kafka.Kafka
}

func (a *API) IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("frontend/views/index.html")
	if err != nil {
		log.Println("Failed to parse template: ", err)
		http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
		return

	}
	if err := t.Execute(w, nil); err != nil {
		log.Println("Failed to parse template: ", err)
		http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
		return
	}
}

func (a *API) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Tried to access post handler with method: ", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Println("Failed to parse form: ", err)
		http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
		return
	}
	message := r.FormValue("message")
	if message == "" {
		log.Println("Message is empty")
		http.Error(w, "Message is empty", http.StatusBadRequest)
		return
	}
	selectedChannels := r.Form["channels"]
	if len(selectedChannels) == 0 {
		log.Println("No channels selected")
		http.Error(w, "No channels selected", http.StatusBadRequest)
		return

	}
	var channels []channel.Channel
	for _, ch := range selectedChannels {
		switch ch {
		case "sms":
			channels = append(channels, channel.SMS)
		case "email":
			channels = append(channels, channel.Email)
		case "slack":
			channels = append(channels, channel.Slack)
		default:
			log.Println("Invalid channel: ", ch)
		}
	}

	kafkaMsg := model.NotificationRequest{
		Message:  message,
		Channels: channels,
	}

	if err := a.Kafka.Produce(kafkaMsg); err != nil {
		log.Println("Failed to produce message: ", err)
		http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
		return
	}

	// Redirect or respond after processing
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
