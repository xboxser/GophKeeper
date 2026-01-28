package model

import "errors"

var (
	ErrCardDuplicate = errors.New("error card duplicate record")
	ErrCardNotFound  = errors.New("error not found card")
)

type CardAPI struct {
	ID         int
	Title      string `json:"title" validate:"required,min=1"`
	Last4      string `json:"last4" validate:"required,min=4,max=4"`
	Number     []byte `json:"number" validate:"required,min=1"`
	Expiry     []byte `json:"expiry" validate:"required,min=1"`
	CVV        []byte `json:"cvv" validate:"required,min=1"`
	CardHolder []byte `json:"card_holder" validate:"required,min=1"`
}

type Card struct {
	ID         int
	Title      string `json:"title" validate:"required"`
	Number     string `json:"number" validate:"required"`
	Expiry     string `json:"expiry" validate:"required"`
	CVV        string `json:"cvv" validate:"required"`
	CardHolder string `json:"card_holder" validate:"required"`
}
