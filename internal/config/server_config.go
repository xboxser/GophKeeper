package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
}

func NewServerConfig() *ServerConfig {
	var cfg ServerConfig
	_ = env.Parse(&cfg)

	clientFlags := flag.NewFlagSet("server", flag.ExitOnError)
	serverAddress := clientFlags.String("a", "localhost:8018", "адрес и порт сервера")

	if cfg.ServerAddress == "" {
		cfg.ServerAddress = *serverAddress
	}

	return &cfg
}
