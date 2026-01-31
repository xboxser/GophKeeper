package repository

import (
	"context"
	"errors"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/db"
)

//go:generate mockgen -source=credential.go -destination=../../../mocks/server/repository/credential_mock.go -package=repository
type CredentialRepository interface {
	AddCredential(ctx context.Context, credential model.CredentialAPI, userID int) error
	DeleteCredential(ctx context.Context, IncID int, userID int) error
	GetCredentials(ctx context.Context, userID int) ([]model.CredentialAPI, error)
	GetCredential(ctx context.Context, IncID int, userID int) (model.CredentialAPI, error)
	UpdateCredential(ctx context.Context, credential model.CredentialAPI, userID int) error
}

type CredentialDB struct {
	DB                db.DB
	CounterRepository CounterRepository
}

func NewCredentialDB(db db.DB, c CounterRepository) *CredentialDB {

	return &CredentialDB{
		DB:                db,
		CounterRepository: c,
	}
}

func (c *CredentialDB) AddCredential(ctx context.Context, credential model.CredentialAPI, userID int) error {
	count, err := c.CounterRepository.GetCounter(ctx, model.CounterCredential, userID)
	if err != nil {
		return errors.Join(err, errors.New("error get counter"))
	}

	query := `INSERT INTO credentials (user_id, login, password, inc_id) VALUES ($1, $2, $3, $4)`
	_, err = c.DB.Exec(ctx, query, userID, credential.Login, credential.Password, count)
	if err != nil {
		return err
	}

	err = c.CounterRepository.IncrementCounter(ctx, model.CounterCredential, userID)
	if err != nil {
		return errors.Join(err, errors.New("error inc counter"))
	}
	return nil
}

func (c *CredentialDB) DeleteCredential(ctx context.Context, IncID int, userID int) error {
	query := `DELETE FROM credentials WHERE inc_id = $1 AND user_id = $2`
	_, err := c.DB.Exec(ctx, query, IncID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (c *CredentialDB) GetCredential(ctx context.Context, IncID int, userID int) (model.CredentialAPI, error) {
	credential := model.CredentialAPI{}
	rows, err := c.DB.Query(ctx, "SELECT login, password FROM credentials WHERE user_id = $1 AND inc_id = $2 LIMIT 1", userID, IncID)
	if err != nil {
		return credential, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&credential.Login, &credential.Password)
		if err != nil {
			return model.CredentialAPI{}, err
		}

	}
	return credential, nil
}

func (c *CredentialDB) GetCredentials(ctx context.Context, userID int) ([]model.CredentialAPI, error) {
	credentials := []model.CredentialAPI{}
	rows, err := c.DB.Query(ctx, "SELECT login, password, inc_id FROM credentials WHERE user_id = $1 ORDER BY id", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		credential := model.CredentialAPI{}
		err := rows.Scan(&credential.Login, &credential.Password, &credential.IncID)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, credential)
	}
	return credentials, nil
}

func (c *CredentialDB) UpdateCredential(ctx context.Context, credential model.CredentialAPI, userID int) error {
	query := `UPDATE credentials SET login = $2, password = $3 WHERE user_id = $1 AND inc_id = $4`
	_, err := c.DB.Exec(ctx, query, userID, credential.Login, credential.Password, credential.IncID)
	if err != nil {
		return err
	}
	return nil
}
