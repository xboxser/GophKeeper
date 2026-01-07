package repository

import (
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/db"
)

type PaymentCardRepository interface {
	GetPaymentCards() ([]model.PaymentCard, error)
	AddPaymentCard(paymentCard model.PaymentCard) error
}

// PaymentCardDB - банковские карты пользователя
type PaymentCardDB struct {
	DB           db.DB
	PaymentCards []model.PaymentCard
}

func NewPaymentCardDB(db db.DB) *PaymentCardDB {
	return &PaymentCardDB{
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

func (p *PaymentCardDB) GetPaymentCards() ([]model.PaymentCard, error) {
	return p.PaymentCards, nil
}

func (p *PaymentCardDB) AddPaymentCard(paymentCard model.PaymentCard) error {
	p.PaymentCards = append(p.PaymentCards, paymentCard)
	return nil
}
