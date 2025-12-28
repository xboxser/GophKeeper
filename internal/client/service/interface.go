package service

import "gophkeeper/internal/model"

type CredentialService interface {
	GetCredentials() ([]model.Credential, error)
}
