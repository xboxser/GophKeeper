package service

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetToken(t *testing.T) {
	client := &http.Client{}
	senderService := NewSenderService(client, "localhost:8080")
	token := "test_token"

	senderService.SetToken(token)

	require.Equal(t, token, senderService.token)
}
