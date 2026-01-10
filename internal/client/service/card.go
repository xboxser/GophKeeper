package service

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/model"
	"net/http"
	"time"
)

type CardService interface {
	GetCards() ([]model.Card, error)
	AddCard(model.Card) error
}

type cardService struct {
	SenderService SenderService
	TokenService  TokenService
}

func NewCardService(senderService SenderService) *cardService {
	return &cardService{
		SenderService: senderService,
		TokenService:  nil,
	}
}

func (s *cardService) InitToken(tokenService TokenService) {
	s.TokenService = tokenService
}

func (s *cardService) GetCards() ([]model.Card, error) {
	token, err := s.TokenService.GetToken()
	if err != nil {
		return nil, err
	}
	s.SenderService.SetToken(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body, response, err := s.SenderService.SendGet(ctx, "/api/card")

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error get card, %v", string(body))
	}

	var cards []model.Card
	if err = json.Unmarshal(body, &cards); err != nil {
		return nil, err
	}

	return cards, nil
}

func (s *cardService) AddCard(card model.Card) error {
	token, err := s.TokenService.GetToken()
	if err != nil {
		return err
	}
	s.SenderService.SetToken(token)

	json, err := json.Marshal(card)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body, response, err := s.SenderService.SendPost(ctx, "/api/card", json)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error add card, %v", string(body))
	}
	return nil
}
