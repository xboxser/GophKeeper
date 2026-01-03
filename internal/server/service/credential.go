package service

import (
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/repository"
)

type CredentialService interface {
	GetCredentials() ([]model.Credential, error)
}

type credentialService struct {
	credentialRepository repository.CredentialRepository
}

func NewCredentialService(credentialRepository repository.CredentialRepository) *credentialService {
	return &credentialService{
		credentialRepository: credentialRepository,
	}
}

func (s *credentialService) GetCredentials() ([]model.Credential, error) {
	return s.credentialRepository.GetCredentials()
}
