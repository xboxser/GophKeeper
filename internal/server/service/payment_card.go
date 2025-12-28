package service

import (
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/repository"
)

type paymentCardService struct {
	PaymentCardRepository repository.PaymentCardRepository
}

func NewPaymentCardService(paymentCardRepository repository.PaymentCardRepository) *paymentCardService {
	return &paymentCardService{
		PaymentCardRepository: paymentCardRepository,
	}
}

func (s *paymentCardService) GetPaymentCards() ([]model.PaymentCard, error) {
	return s.PaymentCardRepository.GetPaymentCards()
}
