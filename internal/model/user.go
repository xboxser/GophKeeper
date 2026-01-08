package model

import "errors"

var (
	ErrRegisterUser   = errors.New("error registering user")
	ErrLoginBusy      = errors.New("login is busy")
	ErrLoginIncorrect = errors.New("incorrect password or login")
)

type User struct {
	ID       int
	Login    string
	Password string
}

type APIUser struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
