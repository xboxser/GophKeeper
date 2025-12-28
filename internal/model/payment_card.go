package model

type PaymentCard struct {
	CardNumber string `json:"card_number"`
	CVV        string `json:"cvv"`
	Expiry     string `json:"expiry"`
	CardHolder string `json:"card_holder"`
}
