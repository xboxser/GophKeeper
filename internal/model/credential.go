package model

import "errors"

var (
	ErrCredentialDuplicate = errors.New("error duplicate record Credential")
	ErrCredentialNotFound  = errors.New("error not found Credential")
	ErrCredentialEmptyID   = errors.New("error empty ID")
)

type CredentialAPI struct {
	IncID    int    `json:"id" validate:""`
	Login    string `json:"login" validate:"required"`
	Password []byte `json:"password" validate:"required"`
}

type Credential struct {
	ID       int
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
