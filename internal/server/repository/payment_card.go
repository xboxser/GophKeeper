package repository

import (
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/db"
)

type PaymentCardRepository interface {
	GetPaymentCards() ([]model.PaymentCard, error)
	AddPaymentCard(paymentCard model.PaymentCard) error
}

// PaymentCardBD - банковские карты пользователя
type PaymentCardBD struct {
	DB           db.DB
	PaymentCards []model.PaymentCard
}

func NewPaymentCardBD(db db.DB) *PaymentCardBD {
	return &PaymentCardBD{
		DB: db,
		PaymentCards: []model.PaymentCard{
			{
				CVV:        "123",
				CardHolder: "Ivan",
				CardNumber: "1234567890123456",
				Expiry:     "01/23",
			},
		},
	}
}

func (p *PaymentCardBD) GetPaymentCards() ([]model.PaymentCard, error) {
	return p.PaymentCards, nil
}

func (p *PaymentCardBD) AddPaymentCard(paymentCard model.PaymentCard) error {
	p.PaymentCards = append(p.PaymentCards, paymentCard)
	return nil
}
