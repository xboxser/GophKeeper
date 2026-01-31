package repository

import (
	"context"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/db"

	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=user.go -destination=../../../mocks/server/repository/user_mock.go -package=repository
type UserRepository interface {
	GetUserForLogin(context.Context, string) (model.User, error)
	RegisterUser(context.Context, model.APIUser) (int, error)
	GetUserForID(context.Context, int) (model.User, error)
}

type UserDB struct {
	DB                db.DB
	CounterRepository CounterRepository
}

func NewUserDB(db db.DB, c CounterRepository) *UserDB {
	return &UserDB{
		DB:                db,
		CounterRepository: c,
	}
}

// RegisterUser - регистрируем пользователя в системе
func (u *UserDB) RegisterUser(ctx context.Context, apiUser model.APIUser) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(apiUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	tx, err := u.DB.Begin(ctx)
	defer tx.Rollback(ctx)

	// Регистрируем пользователя и пытаемся получить его ID
	userID := 0
	query := `INSERT INTO users (login, password, code) VALUES ($1, $2, $3) RETURNING id`

	rows, err := tx.Query(ctx, query, apiUser.Login, string(hashedPassword), apiUser.Code)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return 0, err
		}
	}

	// Если ID равен нулю, то пользователь не зарегистрирован
	if userID == 0 {
		return 0, model.ErrRegisterUser
	}

	// При регистрации пользователя, создаем счетчики
	err = u.CounterRepository.IncrementCounter(tx, ctx, model.CounterCard, userID)
	if err != nil {
		return 0, err
	}
	err = u.CounterRepository.IncrementCounter(tx, ctx, model.CounterCredential, userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// GetUserForLogin - получить пользователя по его логину
func (u *UserDB) GetUserForLogin(ctx context.Context, login string) (model.User, error) {
	var user model.User
	rows, err := u.DB.Query(ctx, "SELECT id, login, password, code FROM users WHERE login = $1", login)
	if err != nil {
		return model.User{}, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&user.ID, &user.Login, &user.Password, &user.Code)
		if err != nil {
			return model.User{}, err
		}
	}

	return user, nil
}

// GetUserForID - получить пользователя по его ID
func (u *UserDB) GetUserForID(ctx context.Context, ID int) (model.User, error) {
	query := `SELECT id, login, password, code FROM users WHERE id = $1 LIMIT 1`
	rows, err := u.DB.Query(ctx, query, ID)
	if err != nil {
		return model.User{}, err
	}
	defer rows.Close()

	var user model.User
	if rows.Next() {
		err := rows.Scan(&user.ID, &user.Login, &user.Password, &user.Code)
		if err != nil {
			return model.User{}, err
		}
	}

	return user, nil
}
