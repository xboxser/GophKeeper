package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gophkeeper/internal/model"
	"gophkeeper/internal/service"
	"net/http"
	"time"
)

type UserService interface {
	Register(login, password, masterPass string) (string, error)
	Login(login, password string) (string, error)
	// Master - метод для проверки валидности мастер пароля
	Master(masterPass string) error
	InitToken(TokenService)
}

// UserMasterService - интерфейс для проверки мастер пароля
type UserMasterService interface {
	// Master - метод для проверки валидности мастер пароля
	Master(masterPass string) error
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

func (s *userService) Register(login, password, masterPass string) (string, error) {
	hashedMasterPass, err := service.ConvertMasterPassToHash(masterPass)
	if err != nil {
		return "", err
	}

	user := model.APIUser{
		Login:    login,
		Password: password,
		Code:     hashedMasterPass,
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

func (s *userService) Master(masterPass string) error {
	token, err := s.TokenService.GetToken()
	if err != nil {
		return err
	}
	s.SenderService.SetToken(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body, response, err := s.SenderService.SendGet(ctx, "/api/user/code")

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error check masterPass, %v", string(body))
	}

	var code model.Code
	if err = json.Unmarshal(body, &code); err != nil {
		return err
	}

	fmt.Println("code", code)
	fmt.Println("body", string(body))

	if err = service.ValidateHash(code.Hash, masterPass); err != nil {
		return err
	}

	return nil
}
