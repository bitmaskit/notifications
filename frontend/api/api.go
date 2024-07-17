package api

import (
	"net/http"

	"github.com/bitmaskit/notifications/frontend/config"
)

type api struct {
	Config *config.FrontendConfig
}

func New(config *config.FrontendConfig) Api {
	return &api{Config: config}
}

type Api interface {
	IndexHandler(w http.ResponseWriter, r *http.Request)
	PostHandler(w http.ResponseWriter, r *http.Request)
}
