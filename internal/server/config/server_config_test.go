package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSeverConfig(t *testing.T) {

	oldArgs := os.Args
	oldAddress := os.Getenv("SERVER_ADDRESS")

	// сохраняем старые значения, чтобы потом восстановить после прохождения тестов
	defer func() {
		os.Args = oldArgs
		os.Setenv("SERVER_ADDRESS", oldAddress)
	}()

	// Проверяем на получение дефолтных значений
	t.Run("checking the receipt of default parameters", func(t *testing.T) {
		os.Unsetenv("SERVER_ADDRESS")
		// Устанавливаем значения, чтобы не было ошибок у парсера флагов
		os.Args = []string{"program"}

		cfg := NewServerConfig()
		require.NotEmpty(t, cfg.ServerAddress)
	})

	// Тестируем, что при наличии переменных окружения, значения они корректно попадают в конфиг
	t.Run("checking set env parameters", func(t *testing.T) {

		URL := "localhost:8080"

		os.Setenv("SERVER_ADDRESS", URL)

		// Устанавливаем значения, чтобы не было ошибок у парсера флагов
		os.Args = []string{"program"}

		cfg := NewServerConfig()
		require.Equal(t, cfg.ServerAddress, URL)
	})
}
