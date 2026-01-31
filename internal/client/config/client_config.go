package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type ClientConfig struct {
	ServerAddress    string `env:"SERVER_ADDRESS"`
	LocalStorage     string `env:"LOCAL_STORAGE"`
	TLS_certFile     string `env:"TLS_CERT_FILE"`
	LimitConTimeFile int    `env:"LIMIT__CON_TIME_FILE"`
}

func NewClientConfig() *ClientConfig {
	var cfg ClientConfig
	_ = env.Parse(&cfg)

	clientFlags := flag.NewFlagSet("client", flag.ExitOnError)
	serverAddress := clientFlags.String("a", "https://localhost:8018", "адрес и порт сервера")
	tlsCertFile := clientFlags.String("tlc_cert", "../../certificate/cert.pem", "Путь до сертификата tlc")
	localStorage := clientFlags.String("ls", "../../upload_client/", "адрес локального файлового хранилища для скачанных файлов с сервера")
	LimitConTimeFile := clientFlags.Int("lct", 1000, "Максимальное время подключение на скачивание/загрузки файла")

	if cfg.ServerAddress == "" {
		cfg.ServerAddress = *serverAddress
	}

	if cfg.LocalStorage == "" {
		cfg.LocalStorage = *localStorage
	}

	if cfg.TLS_certFile == "" {
		cfg.TLS_certFile = *tlsCertFile
	}

	if cfg.LimitConTimeFile == 0 {
		cfg.LimitConTimeFile = *LimitConTimeFile
	}

	return &cfg
}
