package service

import (
	"context"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/repository"
)

type CredentialService interface {
	GetCredentials(ctx context.Context, userID int) ([]model.CredentialAPI, error)
	AddCredential(ctx context.Context, credential model.CredentialAPI, userID int) error
}

type credentialService struct {
	credentialRepository repository.CredentialRepository
}

func NewCredentialService(credentialRepository repository.CredentialRepository) *credentialService {
	return &credentialService{
		credentialRepository: credentialRepository,
	}
}

func (s *credentialService) GetCredentials(ctx context.Context, userID int) ([]model.CredentialAPI, error) {
	return s.credentialRepository.GetCredentials(ctx, userID)
}

func (s *credentialService) AddCredential(ctx context.Context, credential model.CredentialAPI, userID int) error {
	return s.credentialRepository.AddCredential(ctx, credential, userID)

}
