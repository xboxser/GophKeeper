package model

type Card struct {
	ID         int
	Title      string `json:"title" validate:"required"`
	Number     string `json:"number" validate:"required"`
	Expiry     string `json:"expiry" validate:"required"`
	CVV        string `json:"cvv" validate:"required"`
	CardHolder string `json:"card_holder" validate:"required"`
}
