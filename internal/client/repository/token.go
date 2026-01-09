package repository

import (
	"fmt"
	"os"
)

type TokenRepository interface {
	GetToken() (string, error)
	SetToken(string) error
}

type tokenRepository struct {
}

func NewTokenRepository() *tokenRepository {
	return &tokenRepository{}
}

func (t *tokenRepository) GetToken() (string, error) {
	data, err := os.ReadFile("token")
	if err != nil {
		return "", fmt.Errorf("Не удалось прочитать токен: %v", err)
	}
	token := string(data)

	if token == "" {
		return "", fmt.Errorf("Пустой  токен, требуется авторизация")
	}

	return token, nil
}

// SetToken - сохраняем токен пользователя, файл каждый раз перезаписывается
func (t *tokenRepository) SetToken(token string) error {
	err := os.WriteFile("token", []byte(token), 0600) // 0600 только владелец может читать и писать
	if err != nil {
		return fmt.Errorf("Не удалось записать токен в файл: %v", err)
	}
	return nil

}
