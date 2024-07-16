package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

const env = ".env"

var ErrFrontendPortNotSet = errors.New("FRONTEND_PORT is not set")
var ErrBackendHostNotSet = errors.New("BACKEND_HOST is not set")

type FrontendConfig struct {
	Port        string
	BackendAddr string
}

func Load() (*FrontendConfig, error) {
	if err := godotenv.Load(env); err != nil {
		return nil, err
	}

	port := os.Getenv("FRONTEND_PORT")
	if port == "" {
		return nil, ErrFrontendPortNotSet
	}
	backendHost := os.Getenv("BACKEND_HOST")
	if backendHost == "" {
		return nil, ErrBackendHostNotSet
	}
	return &FrontendConfig{
		Port:        port,
		BackendAddr: backendHost,
	}, nil
}
