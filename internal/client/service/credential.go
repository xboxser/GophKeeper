package service

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/client/repository"
	"gophkeeper/internal/model"
	"net/http"
	"time"
)

type CredentialService interface {
	GetCredentials() ([]model.Credential, error)
	AddCredential(login, password string) error
	InitToken(TokenService)
}

type credentialService struct {
	CredentialRepository repository.CredentialRepository
	SenderService        SenderService
	TokenService         TokenService
}

func NewCredentialService(credentialRepository repository.CredentialRepository, senderService SenderService) *credentialService {
	return &credentialService{
		CredentialRepository: credentialRepository,
		SenderService:        senderService,
		TokenService:         nil,
	}
}

// TODO подумать как избавиться от множественного реализации данной функции
// InitToken - добавляем сервис токена для работы с ним
func (s *credentialService) InitToken(tokenService TokenService) {
	s.TokenService = tokenService
}

func (s *credentialService) GetCredentials() ([]model.Credential, error) {
	token, err := s.TokenService.GetToken()
	if err != nil {
		return nil, err
	}
	s.SenderService.SetToken(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body, response, err := s.SenderService.SendGet(ctx, "/api/credentials")

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error get credentials, %v", string(body))
	}

	var credentials []model.Credential
	if err = json.Unmarshal(body, &credentials); err != nil {
		return nil, err
	}

	return credentials, nil
}

func (s *credentialService) AddCredential(login, password string) error {
	token, err := s.TokenService.GetToken()
	if err != nil {
		return err
	}
	s.SenderService.SetToken(token)

	json, err := json.Marshal(model.Credential{Login: login, Password: password})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body, response, err := s.SenderService.SendPost(ctx, "/api/credentials", json)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error add credentials, %v", string(body))
	}
	return nil
}
