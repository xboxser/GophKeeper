package repository

import (
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/db"
)

type CredentialRepository interface {
	GetCredentials() ([]model.Credential, error)
	AddCredential(model.Credential) error
}

type CredentialDB struct {
	DB          db.DB
	Credentials []model.Credential
}

func NewCredentialDB(db db.DB) *CredentialDB {

	return &CredentialDB{
		DB: db,
		Credentials: []model.Credential{
			{Login: "test", Password: "test"},
		},
	}
}

func (c *CredentialDB) GetCredentials() ([]model.Credential, error) {
	return c.Credentials, nil
}

func (c *CredentialDB) AddCredential(credential model.Credential) error {
	c.Credentials = append(c.Credentials, credential)
	return nil
}
