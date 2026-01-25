package service

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/model"
	"net/http"
	"regexp"
	"time"
)

type CardService interface {
	GetCards() ([]model.Card, error)
	AddCard(model.Card) error
	SetMasterPass(string)
}

type cardService struct {
	SenderService     SenderService
	TokenService      TokenService
	EncryptionService EncryptionService
}

func NewCardService(senderService SenderService, encryptionService EncryptionService) *cardService {
	return &cardService{
		SenderService:     senderService,
		EncryptionService: encryptionService,
		TokenService:      nil,
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

	var cardsAPI []model.CardAPI
	if err = json.Unmarshal(body, &cardsAPI); err != nil {
		return nil, err
	}

	var cards []model.Card
	for _, cardAPI := range cardsAPI {
		card, err := s.ConvertToCard(cardAPI)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func (s *cardService) AddCard(card model.Card) error {
	token, err := s.TokenService.GetToken()
	if err != nil {
		return err
	}
	s.SenderService.SetToken(token)

	cardAPI, err := s.ConvertToCardAPI(card)
	if err != nil {
		return err
	}

	json, err := json.Marshal(cardAPI)
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

func (s *cardService) SetMasterPass(masterPass string) {
	s.EncryptionService.SetMasterPass(masterPass)
}

// ConvertToCardAPI - преобразует model.Card в model.CardAPI
// шифрует необходимые поля
func (s *cardService) ConvertToCardAPI(card model.Card) (model.CardAPI, error) {
	cardAPI := model.CardAPI{
		Title: card.Title,
	}
	var err error

	cardAPI.Last4 = getLast4(card.Number)

	cardAPI.Number, err = s.EncryptionService.Encrypt(card.Number)
	if err != nil {
		return model.CardAPI{}, err
	}

	cardAPI.CVV, err = s.EncryptionService.Encrypt(card.CVV)
	if err != nil {
		return model.CardAPI{}, err
	}

	cardAPI.Expiry, err = s.EncryptionService.Encrypt(card.Expiry)
	if err != nil {
		return model.CardAPI{}, err
	}

	cardAPI.CardHolder, err = s.EncryptionService.Encrypt(card.CardHolder)
	if err != nil {
		return model.CardAPI{}, err
	}

	return cardAPI, nil
}

// ConvertToCard - преобразует model.CardAPI в model.Card
// дешифрует нужные поля
func (s *cardService) ConvertToCard(cardAPI model.CardAPI) (model.Card, error) {
	card := model.Card{
		Title: cardAPI.Title,
	}
	var err error

	card.Number, err = s.EncryptionService.Decrypt(cardAPI.Number)
	if err != nil {
		return model.Card{}, err
	}

	card.CVV, err = s.EncryptionService.Decrypt(cardAPI.CVV)
	if err != nil {
		return model.Card{}, err
	}

	card.Expiry, err = s.EncryptionService.Decrypt(cardAPI.Expiry)
	if err != nil {
		return model.Card{}, err
	}

	card.CardHolder, err = s.EncryptionService.Decrypt(cardAPI.CardHolder)
	if err != nil {
		return model.Card{}, err
	}

	return card, nil

}

// getLast4 - получить последние 4 цифры карты
func getLast4(number string) string {
	var digitRegexp = regexp.MustCompile(`[^\d]`)
	number = digitRegexp.ReplaceAllString(number, "")
	return number[len(number)-4:]
}
