package repository

import (
	"context"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/db"
)

type CardRepository interface {
	GetCards(ctx context.Context, userID int) ([]model.Card, error)
	AddCard(ctx context.Context, card model.Card, userID int) error
}

type CardDB struct {
	DB db.DB
}

func NewCardDB(db db.DB) *CardDB {
	return &CardDB{DB: db}
}

func (c *CardDB) GetCards(ctx context.Context, userID int) ([]model.Card, error) {
	var cards []model.Card
	rows, err := c.DB.Query(ctx, "SELECT id, title, number_enc, expiry_enc, cvv_enc, card_holder_name FROM bank_cards WHERE user_id = $1 ORDER BY id", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		card := model.Card{}
		err := rows.Scan(&card.ID, &card.Title, &card.Number, &card.Expiry, &card.CVV, &card.CardHolder)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func (c *CardDB) AddCard(ctx context.Context, card model.Card, userID int) error {
	query := `INSERT INTO bank_cards (user_id, title, number_enc, expiry_enc, cvv_enc, card_holder_name) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := c.DB.Exec(ctx, query, userID, card.Title, card.Number, card.Expiry, card.CVV, card.CardHolder)
	if err != nil {
		return err
	}
	return nil
}
