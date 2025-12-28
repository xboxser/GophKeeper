package repository

import "gophkeeper/internal/model"

type CredentialRepository interface {
	GetCredentials() ([]model.Credential, error)
	AddCredential(credential model.Credential) error
}
