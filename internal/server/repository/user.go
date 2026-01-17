package repository

import (
	"context"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/db"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetUserForLogin(context.Context, string) (model.User, error)
	RegisterUser(context.Context, model.APIUser) (int, error)
	GetUserForID(context.Context, int) (model.User, error)
}

type UserDB struct {
	DB db.DB
}

func NewUserDB(db db.DB) *UserDB {
	return &UserDB{DB: db}
}

// RegisterUser - регистрируем пользователя в системе
func (u *UserDB) RegisterUser(ctx context.Context, apiUser model.APIUser) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(apiUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// Регистрируем пользователя и пытаемся получить его ID
	userID := 0
	query := `INSERT INTO users (login, password, code) VALUES ($1, $2, $3) RETURNING id`

	rows, err := u.DB.Query(ctx, query, apiUser.Login, string(hashedPassword), apiUser.Code)
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
