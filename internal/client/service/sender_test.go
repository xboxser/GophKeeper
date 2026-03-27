package service

import (
	"bytes"
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

func TestSendPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result": "post successful"}`))
	}))
	defer server.Close()

	sender := &senderService{
		client:        &http.Client{Timeout: 5 * time.Second},
		serverAddress: server.URL,
		token:         "Bearer test-token",
	}

	ctx := context.Background()
	url := "/test-post-endpoint"
	data := []byte(`{"key": "value"}`)

	_, response, err := sender.SendPost(ctx, url, data)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestSendPut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что метод запроса PUT
		assert.Equal(t, http.MethodPut, r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result": "put successful"}`))
	}))
	defer server.Close()

	sender := &senderService{
		client:        &http.Client{Timeout: 5 * time.Second},
		serverAddress: server.URL,
		token:         "Bearer test-token",
	}

	ctx := context.Background()
	url := "/test-put-endpoint"
	data := []byte(`{"key": "updated_value"}`)

	_, response, err := sender.SendPut(ctx, url, data)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestSendGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result": "get successful"}`))
	}))
	defer server.Close()

	sender := &senderService{
		client:        &http.Client{Timeout: 5 * time.Second},
		serverAddress: server.URL,
		token:         "Bearer test-token",
	}

	ctx := context.Background()
	url := "/test-get-endpoint"

	_, response, err := sender.SendGet(ctx, url)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestSendFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("file upload successful"))
	}))
	defer server.Close()

	sender := &senderService{
		client:        &http.Client{Timeout: 5 * time.Second},
		serverAddress: server.URL,
		token:         "Bearer test-token",
	}

	ctx := context.Background()
	url := "/test-file-endpoint"
	fileContent := []byte("test file content")
	fileReader := bytes.NewReader(fileContent)
	filename := "test_file.txt"

	_, response, err := sender.SendFile(ctx, url, fileReader, filename)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestSendGetFile(t *testing.T) {
	expectedFileContent := []byte("test file content")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(expectedFileContent)
	}))
	defer server.Close()

	sender := &senderService{
		client:        &http.Client{Timeout: 5 * time.Second},
		serverAddress: server.URL,
		token:         "Bearer test-token",
	}

	ctx := context.Background()
	url := "/test-get-file-endpoint"

	response, err := sender.SendGetFile(ctx, url)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}
