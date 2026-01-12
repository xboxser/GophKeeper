package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTokenRepository(t *testing.T) {

	fileName := "file_test"
	tokenRep := NewTokenRepository(fileName)

	require.Equal(t, fileName, tokenRep.FileName)
}
