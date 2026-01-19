package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type ClientConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	LocalStorage  string `env:"LOCAL_STORAGE"`
}

func NewClientConfig() *ClientConfig {
	var cfg ClientConfig
	_ = env.Parse(&cfg)

	clientFlags := flag.NewFlagSet("client", flag.ExitOnError)
	serverAddress := clientFlags.String("a", "http://localhost:8018", "адрес и порт сервера")
	LocalStorage := clientFlags.String("ls", "../../upload_client/", "адрес локального файлового хранилища для скачанных файлов с сервера")

	if cfg.ServerAddress == "" {
		cfg.ServerAddress = *serverAddress
	}

	if cfg.LocalStorage == "" {
		cfg.LocalStorage = *LocalStorage
	}
	return &cfg
}
