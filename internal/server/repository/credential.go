package repository

import (
	"context"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/db"
)

type CredentialRepository interface {
	GetCredentials(ctx context.Context, userID int) ([]model.Credential, error)
	AddCredential(ctx context.Context, credential model.Credential, userID int) error
}

type CredentialDB struct {
	DB db.DB
}

func NewCredentialDB(db db.DB) *CredentialDB {

	return &CredentialDB{
		DB: db,
	}
}

func (c *CredentialDB) GetCredentials(ctx context.Context, userID int) ([]model.Credential, error) {
	credentials := []model.Credential{}
	rows, err := c.DB.Query(ctx, "SELECT login, password FROM credentials WHERE user_id = $1 ORDER BY id", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		credential := model.Credential{}
		err := rows.Scan(&credential.Login, &credential.Password)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, credential)
	}
	return credentials, nil
}

func (c *CredentialDB) AddCredential(ctx context.Context, credential model.Credential, userID int) error {
	query := `INSERT INTO credentials (user_id, login, password) VALUES ($1, $2, $3)`
	_, err := c.DB.Exec(ctx, query, userID, credential.Login, credential.Password)
	if err != nil {
		return err
	}
	return nil
}
