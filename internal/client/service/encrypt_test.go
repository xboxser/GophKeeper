package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewEncryptionService(t *testing.T) {
	encryptionService := NewEncryptionService()

	require.Equal(t, encryptionService.saltLen, 16)
	require.Equal(t, encryptionService.nonceLen, 12)
}

func TestSetMasterPass(t *testing.T) {
	encryptionService := NewEncryptionService()
	encryptionService.SetMasterPass("test")

	require.Equal(t, encryptionService.masterPass, "test")
}

func TestEncryptDecrypt(t *testing.T) {
	encryptionService := NewEncryptionService()

	text := "text test encrypt"
	encryptedFields, err := encryptionService.Encrypt(text)
	fmt.Println(encryptedFields, err)
	require.Error(t, err)
	require.True(t, len(encryptedFields) == 0)

	encryptionService.SetMasterPass("test")
	encryptedFields, err = encryptionService.Encrypt(text)
	require.NoError(t, err)

	decryptedText, err := encryptionService.Decrypt(encryptedFields)
	require.NoError(t, err)

	require.Equal(t, text, decryptedText)

	encryptionService.SetMasterPass("test2")
	decryptedText, err = encryptionService.Decrypt(encryptedFields)
	require.Error(t, err)
	require.Equal(t, "", decryptedText)
}
