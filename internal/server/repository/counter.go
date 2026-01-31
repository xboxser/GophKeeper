package repository

import (
	"context"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/db"
)

type CounterRepository interface {
	GetCounter(ctx context.Context, typeCounter model.TypeCounter, userID int) (int, error)
	IncrementCounter(ctx context.Context, typeCounter model.TypeCounter, userID int) error
}

type counterDB struct {
	DB db.DB
}

func NewUserCounterDB(db db.DB) *counterDB {
	return &counterDB{DB: db}
}

func (c *counterDB) GetCounter(ctx context.Context, typeCounter model.TypeCounter, userID int) (int, error) {
	var count int
	query := "SELECT count FROM counters WHERE user_id = $1 AND type = $2 LIMIT 1"
	//TODO добавить составной индекс  для user_id и type
	//TODO добавить транзакции
	//TODO добавить блокировку строки
	rows, err := c.DB.Query(ctx, query,
		userID, typeCounter)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}

	}
	return count, nil
}

func (uc *counterDB) IncrementCounter(ctx context.Context, typeCounter model.TypeCounter, userID int) error {
	query := `INSERT INTO counters (user_id, type, count)
		VALUES ($1, $2, 1)
		ON CONFLICT (user_id, type)
		DO UPDATE SET count = counters.count + 1
		RETURNING count;`
	_, err := uc.DB.Exec(ctx, query, userID, typeCounter)
	if err != nil {
		return err
	}
	return nil
}
