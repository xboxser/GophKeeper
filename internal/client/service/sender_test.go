package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetToken(t *testing.T) {
	senderService := NewSenderService("", "localhost:8080")
	token := "test_token"

	senderService.SetToken(token)

	require.Equal(t, token, senderService.token)
}
