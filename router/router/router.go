package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"

	"github.com/bitmaskit/notifications/channel"
	"github.com/bitmaskit/notifications/kafka"
	"github.com/bitmaskit/notifications/model"
	"github.com/bitmaskit/notifications/router/config"
)

func route(kafka kafka.Kafka, msg []byte, cfg *config.RouterConfig) error {
	message, channels, err := decodeMessage(msg)
	if err != nil {
		log.Println("Error decoding message: ", err)
		return err
	}
	var errs []error
	for _, ch := range channels {
		var err error
		switch ch {
		case channel.Email:
			err = kafka.ProduceToTopic(message, cfg.EmailTopic)
		case channel.Slack:
			err = kafka.ProduceToTopic(message, cfg.SlackTopic)
		case channel.SMS:
			err = kafka.ProduceToTopic(message, cfg.SmsTopic)
		default:
			err = errors.New("unknown channel")
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		log.Print("Errors routing message: \n", errs)
		for _, e := range errs {
			log.Println(e)
		}
		return errs[0]
	}
	return nil
}

func decodeMessage(msgValue []byte) (string, []channel.Channel, error) {
	nr := model.NotificationRequest{}
	reader := bytes.NewReader(msgValue)
	if err := json.NewDecoder(reader).Decode(&nr); err != nil {
		log.Println("Error decoding request body:", err)
		return "", nil, err
	}

	return nr.Message, nr.Channels, nil
}
