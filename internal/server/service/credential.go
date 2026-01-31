package service

import (
	"context"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/repository"
)

type CredentialService interface {
	AddCredential(ctx context.Context, credential model.CredentialAPI, userID int) error
	DeleteCredential(ctx context.Context, IncID int, userID int) error
	GetCredentials(ctx context.Context, userID int) ([]model.CredentialAPI, error)
	UpdateCredential(ctx context.Context, credential model.CredentialAPI, userID int) error
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

func (s *credentialService) DeleteCredential(ctx context.Context, IncID int, userID int) error {
	// ищем уже созданных логинов у пользователя
	credentialBD, err := s.credentialRepository.GetCredential(ctx, IncID, userID)
	if err != nil {
		return err
	}
	if credentialBD.Login == "" {
		return model.ErrCredentialNotFound
	}

	return s.credentialRepository.DeleteCredential(ctx, IncID, userID)
}

func (s *credentialService) UpdateCredential(ctx context.Context, credential model.CredentialAPI, userID int) error {
	// ищем наличие уже созданных логинов у пользователя
	credentialBD, err := s.credentialRepository.GetCredential(ctx, credential.IncID, userID)
	if err != nil {
		return err
	}
	if credentialBD.Login == "" {
		return model.ErrCredentialNotFound
	}

	return s.credentialRepository.UpdateCredential(ctx, credential, userID)
}
