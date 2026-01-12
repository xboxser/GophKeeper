package model

import (
	"errors"
	"time"
)

var (
	ErrSessionNotFound     = errors.New("session not found")
	ErrSessionExpired      = errors.New("session expired")
	ErrSessionGenerateUUID = errors.New("error generate uuid")
)

type Session struct {
	UserID    int
	UUID      string
	ExpiresAt time.Time // время жизни сессии
}
