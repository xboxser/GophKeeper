package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTokenRepository(t *testing.T) {

	fileName := "file_test"
	tokenRep := NewTokenRepository(fileName)

	require.Equal(t, fileName, tokenRep.FileName)
}

func TestGetToken(t *testing.T) {
	fileName := "file_test"

	t.Run("empty token", func(t *testing.T) {
		tokenRep := NewTokenRepository(fileName)
		token := ""

		err := os.WriteFile(fileName, []byte(token), 0600)
		require.Nil(t, err)

		_, err = tokenRep.GetToken()
		require.Error(t, err)
	})

	t.Run("valid token", func(t *testing.T) {
		tokenRep := NewTokenRepository(fileName)
		token := "123__()(*(&*^%^))"

		err := os.WriteFile(fileName, []byte(token), 0600)
		require.Nil(t, err)

		tokenFile, err := tokenRep.GetToken()
		require.Nil(t, err)
		require.Equal(t, token, tokenFile)
	})

	err := os.Remove(fileName)
	require.Nil(t, err)
}

func TestSetToken(t *testing.T) {
	fileName := "file_test"

	tokenRep := NewTokenRepository(fileName)
	token := "123__()(*(&*^%^))"

	err := tokenRep.SetToken(token)
	require.Nil(t, err)

	tokenFile, err := tokenRep.GetToken()
	require.Nil(t, err)
	require.Equal(t, token, tokenFile)

	err = os.Remove(fileName)
	require.Nil(t, err)
}
