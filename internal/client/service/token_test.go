package service

import (
	"gophkeeper/internal/client/repository"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSetToken(t *testing.T) {
	tokenRep := repository.NewTokenRepository("test")
	tokenService := NewTokenService(tokenRep)
	token := "test_token"

	defer os.Remove("test")
	err := tokenService.SetToken(token)
	require.NoError(t, err)

	tokenGet, err := tokenService.GetToken()
	require.NoError(t, err)
	require.Equal(t, tokenGet, token)
}
