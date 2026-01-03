package repository

import "gophkeeper/internal/model"

type CredentialRepository interface {
	GetCredentials() ([]model.Credential, error)
	AddCredential(credential model.Credential) error
}

type CredentialBD struct {
	Credentials []model.Credential
}

func NewCredentialBD() *CredentialBD {
	return &CredentialBD{
		Credentials: []model.Credential{
			{Login: "test", Password: "test"},
		},
	}
}

func (c *CredentialBD) GetCredentials() ([]model.Credential, error) {
	return c.Credentials, nil
}

func (c *CredentialBD) AddCredential(credential model.Credential) error {
	c.Credentials = append(c.Credentials, credential)
	return nil
}
