package repository

import (
	"context"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/db"
)

type CardRepository interface {
	// AddCard - Добавить карту
	AddCard(ctx context.Context, card model.CardAPI, userID int) error
	// GetCard - Получить информацию по карте пользователя по 4 символам
	DeleteCard(ctx context.Context, last4 string, userID int) error
	GetCard(ctx context.Context, last4 string, userID int) (model.CardAPI, error)
	// GetCards - Получить все карты пользователя
	GetCards(ctx context.Context, userID int) ([]model.CardAPI, error)
}

type CardDB struct {
	DB db.DB
}

func NewCardDB(db db.DB) *CardDB {
	return &CardDB{DB: db}
}

func (c *CardDB) AddCard(ctx context.Context, card model.CardAPI, userID int) error {
	query := `INSERT INTO bank_cards (user_id, title, number_enc, expiry_enc, cvv_enc, card_holder_name, last4) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := c.DB.Exec(ctx, query, userID, card.Title, card.Number, card.Expiry, card.CVV, card.CardHolder, card.Last4)
	if err != nil {
		return err
	}
	return nil
}

func (c *CardDB) DeleteCard(ctx context.Context, last4 string, userID int) error {
	query := `DELETE FROM bank_cards WHERE user_id = $1 AND last4 = $2`
	_, err := c.DB.Exec(ctx, query, userID, last4)
	if err != nil {
		return err
	}
	return nil
}

func (c *CardDB) GetCard(ctx context.Context, last4 string, userID int) (model.CardAPI, error) {
	var card model.CardAPI
	rows, err := c.DB.Query(ctx, `SELECT 
		id, title, number_enc, expiry_enc, cvv_enc, card_holder_name, last4 
		FROM bank_cards WHERE user_id = $1 AND last4 = $2
		LIMIT 1`, userID, last4)
	if err != nil {
		return model.CardAPI{}, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&card.ID, &card.Title, &card.Number, &card.Expiry, &card.CVV, &card.CardHolder, &card.Last4)
		if err != nil {
			return model.CardAPI{}, err
		}
	}
	return card, nil
}

func (c *CardDB) GetCards(ctx context.Context, userID int) ([]model.CardAPI, error) {
	var cards []model.CardAPI
	rows, err := c.DB.Query(ctx, "SELECT id, title, number_enc, expiry_enc, cvv_enc, card_holder_name, last4 FROM bank_cards WHERE user_id = $1 ORDER BY id", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		card := model.CardAPI{}
		err := rows.Scan(&card.ID, &card.Title, &card.Number, &card.Expiry, &card.CVV, &card.CardHolder, &card.Last4)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}
