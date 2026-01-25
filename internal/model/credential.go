package model

import "errors"

var (
	ErrCredentialDuplicate  = errors.New("error duplicate record Credential")
	ErrCredentialNotFound   = errors.New("error not found Credential")
	ErrCredentialEmptyLogin = errors.New("error empty login")
)

type CredentialAPI struct {
	Login    string `json:"login" validate:"required"`
	Password []byte `json:"password" validate:"required"`
}

type Credential struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
