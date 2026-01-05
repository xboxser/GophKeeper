package repository

import (
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/db"
)

type CredentialRepository interface {
	GetCredentials() ([]model.Credential, error)
	AddCredential(credential model.Credential) error
}

type CredentialBD struct {
	DB          db.DB
	Credentials []model.Credential
}

func NewCredentialBD(db db.DB) *CredentialBD {

	return &CredentialBD{
		DB: db,
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
