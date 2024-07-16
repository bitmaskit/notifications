package model

import (
	"io"
	"log"

	"encoding/json"
	"github.com/bitmaskit/notifications/channel"
)

type NotificationRequest struct {
	Message  string            `json:"message"`
	Channels []channel.Channel `json:"channels"`
}

func (r *NotificationRequest) ToJSON() ([]byte, error) {
	jsonData, err := json.Marshal(r)
	if err != nil {
		log.Println("Error marshaling request body:", err)
		return []byte{}, err
	}
	return jsonData, nil
}

func (r *NotificationRequest) ToJSONString() (string, error) {
	jsonData, err := r.ToJSON()
	return string(jsonData), err
}

func (r *NotificationRequest) FromJSON(body io.ReadCloser) (NotificationRequest, error) {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(r)
	if err != nil {
		log.Println("Error decoding request body:", err)

	}
	return *r, err
}
