package config

func New(backendAddr string) *Config {
	return &Config{backendAddr}
}

type Config struct {
	BackendAddr string
}
