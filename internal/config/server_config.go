package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	ServerAddress  string        `env:"SERVER_ADDRESS"`
	DSN            string        `env:"DSN"`
	JWT_SECRET     string        `env:"JWT_SECRET"`
	TokenExpiresAt time.Duration `env:"TOKEN_EXPIRES_AT"`
}

func NewServerConfig() *ServerConfig {
	var cfg ServerConfig
	_ = env.Parse(&cfg)

	clientFlags := flag.NewFlagSet("server", flag.ExitOnError)
	serverAddress := clientFlags.String("a", "localhost:8018", "адрес и порт сервера")
	dsn := clientFlags.String("dsn", "postgres://keeper:YB4df*5f1f)8g92@localhost:5432/GoKeeperDB?sslmode=disable", "подключение к базе данных")

	jwtSecret := clientFlags.String("jwt", "super_secret_code_JWT", " секретный ключ для проверки подписи jwt токена")
	tokenExpiresAt := clientFlags.String("exp", "3h", "время жизни токена")

	if cfg.ServerAddress == "" {
		cfg.ServerAddress = *serverAddress
	}

	if cfg.DSN == "" {
		cfg.DSN = *dsn
	}

	if cfg.JWT_SECRET == "" {
		cfg.JWT_SECRET = *jwtSecret
	}

	if cfg.TokenExpiresAt == 0 {
		cfg.TokenExpiresAt, _ = time.ParseDuration(*tokenExpiresAt)
	}
	return &cfg
}
