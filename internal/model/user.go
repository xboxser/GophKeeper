package model

import "errors"

var (
	ErrUserEmptyCode  = errors.New("error empty field code")
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
