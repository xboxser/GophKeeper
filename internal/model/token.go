package model

import "errors"

var (
	ErrTokenValid = errors.New("token is not valid")
)

type TokenAuth struct {
	UserID int
}
