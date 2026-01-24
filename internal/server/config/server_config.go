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
	UploadPath     string        `env:"UPLOAD_PATH"`
	TLS_certFile   string        `env:"TLS_CERT_FILE"`
	TLS_keyFile    string        `env:"TLS_KEY_FILE"`
}

func NewServerConfig() *ServerConfig {
	var cfg ServerConfig
	_ = env.Parse(&cfg)

	clientFlags := flag.NewFlagSet("server", flag.ExitOnError)
	serverAddress := clientFlags.String("a", "localhost:8018", "адрес и порт сервера")
	dsn := clientFlags.String("dsn", "postgres://keeper:YB4df*5f1f)8g92@localhost:5432/GoKeeperDB?sslmode=disable", "подключение к базе данных")

	jwtSecret := clientFlags.String("jwt", "super_secret_code_JWT", " секретный ключ для проверки подписи jwt токена")
	tokenExpiresAt := clientFlags.String("exp", "3h", "время жизни токена")
	uploadPath := clientFlags.String("upload", "../../upload/", "Путь сохранения файлов")

	tlsCertFile := clientFlags.String("tlc_cert", "../../certificate/cert.pem", "Путь до сертификата tlc")
	tlsKeyFile := clientFlags.String("tlc_key", "../../certificate/key.pem", "Путь до ключа сертификата tlc")

	if cfg.ServerAddress == "" {
		cfg.ServerAddress = *serverAddress
	}

	if cfg.DSN == "" {
		cfg.DSN = *dsn
	}

	if cfg.JWT_SECRET == "" {
		cfg.JWT_SECRET = *jwtSecret
	}

	if cfg.UploadPath == "" {
		cfg.UploadPath = *uploadPath
	}

	if cfg.TokenExpiresAt == 0 {
		cfg.TokenExpiresAt, _ = time.ParseDuration(*tokenExpiresAt)
	}

	if cfg.TLS_certFile == "" {
		cfg.TLS_certFile = *tlsCertFile
	}

	if cfg.TLS_keyFile == "" {
		cfg.TLS_keyFile = *tlsKeyFile
	}

	return &cfg
}
