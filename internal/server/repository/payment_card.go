package repository

import "gophkeeper/internal/model"

type PaymentCardRepository interface {
	GetPaymentCards() ([]model.PaymentCard, error)
	AddPaymentCard(paymentCard model.PaymentCard) error
}

// PaymentCardBD - банковские карты пользователя
type PaymentCardBD struct {
	PaymentCards []model.PaymentCard
}

func NewPaymentCardBD() *PaymentCardBD {
	return &PaymentCardBD{
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
