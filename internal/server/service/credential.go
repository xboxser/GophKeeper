package service

import (
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/repository"
)

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
