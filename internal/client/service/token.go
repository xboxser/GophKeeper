package service

import "gophkeeper/internal/client/repository"

//go:generate mockgen -source=token.go -destination=../../../mocks/client/service/token_mock.go -package=service
type TokenService interface {
	GetToken() (string, error)
	SetToken(string) error
}

type tokenService struct {
	TokenRepository repository.TokenRepository
}

func NewTokenService(tokenRepository repository.TokenRepository) *tokenService {
	return &tokenService{
		TokenRepository: tokenRepository,
	}
}

func (t *tokenService) GetToken() (string, error) {
	return t.TokenRepository.GetToken()
}

func (t *tokenService) SetToken(token string) error {
	return t.TokenRepository.SetToken(token)
}
