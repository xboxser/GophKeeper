package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

type EncryptionService interface {
	Encrypt(string) ([]byte, error)
	Decrypt([]byte) (string, error)
}

type encryptionService struct {
	masterPass string
	saltLen    int
	nonceLen   int
}

func NewEncryptionService(masterPass string) *encryptionService {
	return &encryptionService{
		masterPass: masterPass,
		saltLen:    16,
		nonceLen:   12,
	}
}

func (es *encryptionService) Encrypt(text string) ([]byte, error) {

	// создаем соль случайным способом
	salt := make([]byte, es.saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, errors.New("error generate salt")
	}

	key := es.getArgon2Key(salt)

	// AES-256-GCM шифрование (стандарт PCI DSS 3.4)
	// создание AES-шифровальщика (cipher) из заданного ключа.
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("error generate aes")
	}

	// оборачиваем блочный шифр в режим GCM — Galois/Counter Mode для симметричного шифрования
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("error generate aes")
	}

	// случайное значение,  использоваться только один раз
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.New("error generate nonce")
	}

	// шифруем текст
	cipherText := aesgcm.Seal(nil, nonce, []byte(text), nil)

	// Собираем encrypted blob: salt(16B) + nonce(12B) + ciphertext
	encryptedFields := append(salt, append(nonce, cipherText...)...)

	return encryptedFields, nil
}

func (es *encryptionService) Decrypt(encryptedFields []byte) (string, error) {
	if len(encryptedFields) < es.nonceLen+es.saltLen+1 {
		return "", errors.New("encrypted data too short")
	}

	// извлекаем salt (первые 16 байт)
	salt := encryptedFields[:es.saltLen]

	key := es.getArgon2Key(salt)

	//Извлекаем nonce (16-28 байт) и cipherText (с 28-го байта)
	nonce := encryptedFields[es.saltLen : es.saltLen+es.nonceLen]
	cipherText := encryptedFields[es.saltLen+es.nonceLen:]

	// Восстанавливаем AES-GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes init: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("gcm init: %w", err)
	}

	plaintext, err := aesgcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed (wrong password?): %w", err)
	}

	return string(plaintext), nil
}

func (es *encryptionService) getArgon2Key(salt []byte) []byte {
	// получение криптографического ключа
	// time - Количество итераций (time cost)
	// memory - Количество памяти (memory cost)
	return argon2.IDKey([]byte(es.masterPass), salt, 2, 64*1024, 1, 32)
}
