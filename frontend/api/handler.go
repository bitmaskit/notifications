package api

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/bitmaskit/notifications/channel"
	"github.com/bitmaskit/notifications/model"
)

const (
	ErrInternalServerError = "Internal server error"
	PostEndpoint           = "/message" // TODO: move to config
)

func (a *api) IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/index.html")
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

func (a *api) PostHandler(w http.ResponseWriter, r *http.Request) {
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

	requestBody := model.NotificationRequest{
		Message:  message,
		Channels: channelsFromArray(selectedChannels),
	}

	jsonBody, err := requestBody.ToJSON()
	if err != nil {
		log.Println("Failed to marshal request body: ", err)
		http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest(http.MethodPost, a.Config.BackendAddr+PostEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println("Failed to create request: ", err)
		http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send request: ", err)
		http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Println("Failed to send request: ", resp.Status)
		http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func channelsFromArray(chans []string) []channel.Channel {
	var channels []channel.Channel
	for _, ch := range chans {
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
	return channels
}
