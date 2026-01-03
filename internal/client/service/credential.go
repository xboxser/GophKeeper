package service

import (
	"gophkeeper/internal/client/repository"
	"gophkeeper/internal/model"
)

type CredentialService interface {
	GetCredentials() ([]model.Credential, error)
}

type credentialService struct {
	CredentialRepository repository.CredentialRepository
}

func NewCredentialService(credentialRepository repository.CredentialRepository) *credentialService {
	return &credentialService{
		CredentialRepository: credentialRepository,
	}
}

func (s *credentialService) GetCredentials() ([]model.Credential, error) {
	return s.CredentialRepository.GetCredentials()
}

func (s *credentialService) AddCredential(credential model.Credential) error {
	return s.CredentialRepository.AddCredential(credential)
}
