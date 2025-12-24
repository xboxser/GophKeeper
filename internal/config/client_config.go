package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type ClientConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
}

func NewClientConfig() *ClientConfig {
	var cfg ClientConfig
	_ = env.Parse(&cfg)

	clientFlags := flag.NewFlagSet("server", flag.ExitOnError)
	serverAddress := clientFlags.String("a", "localhost:8008", "адрес и порт сервера")

	if cfg.ServerAddress == "" {
		cfg.ServerAddress = *serverAddress
	}

	return &cfg
}
