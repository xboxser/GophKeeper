package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	DSN           string `env:"DSN"`
}

func NewServerConfig() *ServerConfig {
	var cfg ServerConfig
	_ = env.Parse(&cfg)

	clientFlags := flag.NewFlagSet("server", flag.ExitOnError)
	serverAddress := clientFlags.String("a", "localhost:8018", "адрес и порт сервера")
	dsn := clientFlags.String("dsn", "postgres://keeper:YB4df*5f1f)8g92@localhost:5432/GoKeeperDB?sslmode=disable", "подключение к базе данных")

	if cfg.ServerAddress == "" {
		cfg.ServerAddress = *serverAddress
	}

	if cfg.DSN == "" {
		cfg.DSN = *dsn
	}

	return &cfg
}
