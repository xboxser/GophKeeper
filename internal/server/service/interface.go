package service

import "gophkeeper/internal/model"

type CredentialService interface {
	GetCredentials() ([]model.Credential, error)
}

type PaymentCardService interface {
	GetPaymentCards() ([]model.PaymentCard, error)
}
