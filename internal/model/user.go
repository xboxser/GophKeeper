package model

import "errors"

var (
	ErrRegisterUser = errors.New("error registering user")
	ErrLoginBusy    = errors.New("login is busy")
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
