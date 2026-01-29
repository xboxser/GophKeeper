package repository

import (
	"context"
	"errors"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/db"
)

type CardRepository interface {
	// AddCard - Добавить карту
	AddCard(ctx context.Context, card model.CardAPI, userID int) error
	// GetCard - Получить информацию по карте пользователя по 4 символам
	DeleteCard(ctx context.Context, id string, userID int) error
	GetCard(ctx context.Context, id string, userID int) (model.CardAPI, error)
	// GetCards - Получить все карты пользователя
	GetCards(ctx context.Context, userID int) ([]model.CardAPI, error)

	UpdateCard(ctx context.Context, card model.CardAPI, userID int) error
}

type CardDB struct {
	DB                db.DB
	CounterRepository CounterRepository
}

func NewCardDB(db db.DB, c CounterRepository) *CardDB {
	return &CardDB{
		DB:                db,
		CounterRepository: c,
	}
}

func (c *CardDB) AddCard(ctx context.Context, card model.CardAPI, userID int) error {

	count, err := c.CounterRepository.GetCounter(ctx, model.CounterCard, userID)
	if err != nil {
		return errors.Join(err, errors.New("error get counter"))
	}
	query := `INSERT INTO bank_cards (user_id, title, number_enc, expiry_enc, cvv_enc, card_holder_name, last4, inc_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = c.DB.Exec(ctx, query, userID, card.Title, card.Number, card.Expiry, card.CVV, card.CardHolder, card.Last4, count)
	if err != nil {
		return errors.Join(err, errors.New("error add card"))
	}

	err = c.CounterRepository.IncrementCounter(ctx, model.CounterCard, userID)
	if err != nil {
		return errors.Join(err, errors.New("error inc counter"))
	}
	return nil
}

func (c *CardDB) DeleteCard(ctx context.Context, id string, userID int) error {
	query := `DELETE FROM bank_cards WHERE user_id = $1 AND inc_id = $2`
	_, err := c.DB.Exec(ctx, query, userID, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *CardDB) GetCard(ctx context.Context, id string, userID int) (model.CardAPI, error) {
	var card model.CardAPI
	rows, err := c.DB.Query(ctx, `SELECT 
		id, title, number_enc, expiry_enc, cvv_enc, card_holder_name, last4, inc_id 
		FROM bank_cards WHERE user_id = $1 AND inc_id = $2
		LIMIT 1`, userID, id)
	if err != nil {
		return model.CardAPI{}, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&card.ID, &card.Title, &card.Number, &card.Expiry, &card.CVV, &card.CardHolder, &card.Last4, &card.IncID)
		if err != nil {
			return model.CardAPI{}, err
		}
	}
	return card, nil
}

func (c *CardDB) GetCards(ctx context.Context, userID int) ([]model.CardAPI, error) {
	var cards []model.CardAPI
	rows, err := c.DB.Query(ctx, "SELECT id, title, number_enc, expiry_enc, cvv_enc, card_holder_name, last4, inc_id FROM bank_cards WHERE user_id = $1 ORDER BY id", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		card := model.CardAPI{}
		err := rows.Scan(&card.ID, &card.Title, &card.Number, &card.Expiry, &card.CVV, &card.CardHolder, &card.Last4, &card.IncID)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func (c *CardDB) UpdateCard(ctx context.Context, card model.CardAPI, userID int) error {
	query := `UPDATE bank_cards 
	SET user_id = $1, title = $2, number_enc = $3, expiry_enc = $4, cvv_enc = $5, card_holder_name = $6, last4=$7
	WHERE user_id = $1 AND inc_id = $8`
	_, err := c.DB.Exec(ctx, query, userID, card.Title, card.Number, card.Expiry, card.CVV, card.CardHolder, card.Last4, card.IncID)
	if err != nil {
		return err
	}
	return nil
}
