package model

import (
	"testing"

	"github.com/bitmaskit/notifications/channel"
	"github.com/stretchr/testify/assert"
)

var jsonStr = `{"message":"Hello, World!","channels":["email"]}`

func TestNotificationRequest_FromJSON(t *testing.T) {
	r := NotificationRequest{}

	_, err := r.FromJSON(jsonStr)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", r.Message)
	assert.Equal(t, 1, len(r.Channels))
	assert.Equal(t, "email", r.Channels[0].String())
}

func TestNotificationRequest_ToJSON(t *testing.T) {
	r := NotificationRequest{
		Message:  "Hello, World!",
		Channels: []channel.Channel{channel.Email},
	}

	got, err := r.ToJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte(jsonStr), got)
}
func TestNotificationRequest_ToJSONString(t *testing.T) {
	r := NotificationRequest{
		Message:  "Hello, World!",
		Channels: []channel.Channel{channel.Email},
	}

	got, err := r.ToJSONString()
	assert.NoError(t, err)
	assert.Equal(t, jsonStr, got)
}
