package router

import (
	"errors"
	"log"

	"github.com/bitmaskit/notifications/channel"
	"github.com/bitmaskit/notifications/kafka"
	"github.com/bitmaskit/notifications/kafka/topic"
)

const (
	brokerAddress = "localhost:9092"
)

func Route(msg string, channels []channel.Channel) error {
	kafka := kafka.New(brokerAddress)
	var errs []error
	for _, ch := range channels {
		var err error
		switch ch {
		case channel.Email:
			err = kafka.ProduceToTopic(msg, topic.EmailTopic)
		case channel.Slack:
			err = kafka.ProduceToTopic(msg, topic.SlackTopic)
		case channel.SMS:
			err = kafka.ProduceToTopic(msg, topic.SmsTopic)
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
