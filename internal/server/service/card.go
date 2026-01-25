package service

import (
	"context"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/repository"
)

type CardService interface {
	GetCards(ctx context.Context, userID int) ([]model.CardAPI, error)
	AddCard(ctx context.Context, card *model.CardAPI, userID int) error
}

type cardService struct {
	CardRepository repository.CardRepository
}

func NewCardService(cardRepository repository.CardRepository) *cardService {
	return &cardService{
		CardRepository: cardRepository,
	}
}

func (s *cardService) GetCards(ctx context.Context, userID int) ([]model.CardAPI, error) {
	return s.CardRepository.GetCards(ctx, userID)
}

func (s *cardService) AddCard(ctx context.Context, card *model.CardAPI, userID int) error {
	//TODO добавить проверку полей card, убрать лишние символы
	cardDB, err := s.CardRepository.GetCard(ctx, card.Last4, userID)
	if err != nil {
		return err
	}
	if cardDB.ID != 0 {
		return model.ErrCardDuplicate
	}
	return s.CardRepository.AddCard(ctx, *card, userID)
}
