package service

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/model"
	"net/http"
	"strconv"
	"time"
)

type CredentialService interface {
	AddCredential(login, password, masterPass string) error
	DeleteCredential(id string) error
	GetCredentials(string) ([]model.Credential, error)

	UpdateCredential(login, password, masterPass, id string) error
	InitToken(TokenService)
}

type credentialService struct {
	SenderService     SenderService
	TokenService      TokenService
	EncryptionService EncryptionService
}

func NewCredentialService(senderService SenderService, encryptionService EncryptionService) *credentialService {
	return &credentialService{
		SenderService:     senderService,
		EncryptionService: encryptionService,
		TokenService:      nil,
	}
}

// TODO подумать как избавиться от множественного реализации данной функции
// InitToken - добавляем сервис токена для работы с ним
func (s *credentialService) InitToken(tokenService TokenService) {
	s.TokenService = tokenService
}

func (s *credentialService) DeleteCredential(id string) error {
	token, err := s.TokenService.GetToken()
	if err != nil {
		return err
	}
	s.SenderService.SetToken(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body, response, err := s.SenderService.SendDelete(ctx, "/api/credentials/"+id)

	if err != nil {
		return err
	}

	fmt.Println("status code", response.StatusCode)

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error delete credentials, %v", string(body))
	}
	return nil
}

func (s *credentialService) GetCredentials(masterPass string) ([]model.Credential, error) {
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

	var credentialsAPI []model.CredentialAPI
	if err = json.Unmarshal(body, &credentialsAPI); err != nil {
		return nil, err
	}

	var credentials []model.Credential
	s.EncryptionService.SetMasterPass(masterPass)
	for _, credentialAPI := range credentialsAPI {
		password, err := s.EncryptionService.Decrypt(credentialAPI.Password)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, model.Credential{
			ID:       credentialAPI.IncID,
			Login:    credentialAPI.Login,
			Password: password,
		})
	}

	return credentials, nil
}

func (s *credentialService) AddCredential(login, password, masterPass string) error {
	token, err := s.TokenService.GetToken()
	if err != nil {
		return err
	}
	s.SenderService.SetToken(token)

	s.EncryptionService.SetMasterPass(masterPass)
	passHash, err := s.EncryptionService.Encrypt(password)
	if err != nil {
		return err
	}

	json, err := json.Marshal(model.CredentialAPI{Login: login, Password: passHash})
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

func (s *credentialService) UpdateCredential(login, password, masterPass, id string) error {
	token, err := s.TokenService.GetToken()
	if err != nil {
		return err
	}
	s.SenderService.SetToken(token)

	s.EncryptionService.SetMasterPass(masterPass)
	passHash, err := s.EncryptionService.Encrypt(password)
	if err != nil {
		return err
	}

	incID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	json, err := json.Marshal(model.CredentialAPI{Login: login, Password: passHash, IncID: incID})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body, response, err := s.SenderService.SendPut(ctx, "/api/credentials", json)

	if err != nil {
		return err
	}
	fmt.Println("status code", response.StatusCode)

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error update credentials, %v", string(body))
	}
	return nil
}
