package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetToken(t *testing.T) {
	senderService := NewSenderService("", "localhost:8080")
	token := "test_token"

	senderService.SetToken(token)

	require.Equal(t, token, senderService.token)
}

func TestSendDelete(t *testing.T) {
	// Создаем тестовый HTTP-сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что метод запроса DELETE
		assert.Equal(t, http.MethodDelete, r.Method)

		// Устанавливаем статус ответа
		w.WriteHeader(http.StatusNoContent)
		// Отправляем тело ответа
		_, _ = w.Write([]byte("delete successful"))
	}))
	defer server.Close()

	// Создаем экземпляр senderService
	sender := &senderService{
		client:        &http.Client{Timeout: 5 * time.Second},
		serverAddress: server.URL,
		token:         "Bearer test-token",
	}

	ctx := context.Background()
	url := "/test-endpoint"

	_, response, err := sender.SendDelete(ctx, url)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, response.StatusCode)
}
