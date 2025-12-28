package repository

import "gophkeeper/internal/model"

type CredentialMem struct {
	Credentials []model.Credential
}

func NewCredentialMem() *CredentialMem {
	return &CredentialMem{
		Credentials: []model.Credential{
			{Login: "test", Password: "test"},
		},
	}
}

func (c *CredentialMem) GetCredentials() ([]model.Credential, error) {
	return c.Credentials, nil
}

func (c *CredentialMem) AddCredential(credential model.Credential) error {
	c.Credentials = append(c.Credentials, credential)
	return nil
}
