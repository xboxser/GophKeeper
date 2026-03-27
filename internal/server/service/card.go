package service

import (
	"context"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/repository"
	"strconv"
)

//go:generate mockgen -source=card.go -destination=../../../mocks/server/service/card_mock.go -package=service
type CardService interface {
	AddCard(ctx context.Context, card *model.CardAPI, userID int) error
	DeleteCard(ctx context.Context, id string, userID int) error
	GetCards(ctx context.Context, userID int) ([]model.CardAPI, error)
	UpdateCard(ctx context.Context, card *model.CardAPI, userID int) error
}

type cardService struct {
	CardRepository repository.CardRepository
}

func NewCardService(cardRepository repository.CardRepository) *cardService {
	return &cardService{
		CardRepository: cardRepository,
	}
}

func (s *cardService) AddCard(ctx context.Context, card *model.CardAPI, userID int) error {
	return s.CardRepository.AddCard(ctx, *card, userID)
}

func (s *cardService) DeleteCard(ctx context.Context, id string, userID int) error {
	cardDB, err := s.CardRepository.GetCard(ctx, id, userID)
	if err != nil {
		return err
	}
	if cardDB.ID == 0 {
		return model.ErrCardNotFound
	}
	return s.CardRepository.DeleteCard(ctx, id, userID)
}

func (s *cardService) GetCards(ctx context.Context, userID int) ([]model.CardAPI, error) {
	return s.CardRepository.GetCards(ctx, userID)
}

func (s *cardService) UpdateCard(ctx context.Context, card *model.CardAPI, userID int) error {
	id := strconv.Itoa(card.IncID)
	cardDB, err := s.CardRepository.GetCard(ctx, id, userID)
	if err != nil {
		return err
	}

	if cardDB.ID == 0 {
		return model.ErrCardNotFound
	}
	return s.CardRepository.UpdateCard(ctx, *card, userID)
}
