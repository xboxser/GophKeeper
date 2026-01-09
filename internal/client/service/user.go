package service

import (
	"context"
	"encoding/json"
	"errors"
	"gophkeeper/internal/model"
	"net/http"
	"time"
)

type UserService interface {
	Register(login, password string) (string, error)
	Login(login, password string) (string, error)
	InitToken(TokenService)
}

type userService struct {
	SenderService SenderService
	TokenService  TokenService
}

func NewUserService(sender SenderService) *userService {
	return &userService{
		SenderService: sender,
		TokenService:  nil,
	}
}

// InitToken - добавляем сервис токена для работы с ним
func (s *userService) InitToken(tokenService TokenService) {
	s.TokenService = tokenService
}

func (s *userService) Register(login, password string) (string, error) {
	user := model.APIUser{
		Login:    login,
		Password: password,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	json, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	_, response, err := s.SenderService.SendPost(ctx, "/api/user/register", json)
	if err != nil {
		return "", err
	}

	token := response.Header.Get("Authorization")
	if token == "" {
		return "", errors.New("empty token service")
	}

	if s.TokenService != nil {
		s.TokenService.SetToken(token)
	}

	return token, nil
}

func (s *userService) Login(login, password string) (string, error) {
	user := model.APIUser{
		Login:    login,
		Password: password,
	}

	json, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, response, err := s.SenderService.SendPost(ctx, "/api/user/login", json)
	if err != nil {
		return "", err
	}

	token := response.Header.Get("Authorization")

	if response.StatusCode == http.StatusUnauthorized {
		return "", model.ErrLoginIncorrect
	}

	if token == "" {
		return "", errors.New("empty token service")
	}

	if s.TokenService != nil {
		s.TokenService.SetToken(token)
	}

	return token, nil

}
