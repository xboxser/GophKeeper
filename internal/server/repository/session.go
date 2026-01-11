package repository

import (
	"context"
	"gophkeeper/internal/model"
	"time"
)

type SessionRepository interface {
	GetSession(ctx context.Context, code string) (model.Session, error)
	AddSession(ctx context.Context, session model.Session) error
	DeleteSession(ctx context.Context, code string) error
}

type SessionMemory struct {
	sessions map[string]model.Session
}

func NewSessionMemory() *SessionMemory {
	return &SessionMemory{
		sessions: make(map[string]model.Session),
	}
}

func (sm *SessionMemory) GetSession(ctx context.Context, code string) (model.Session, error) {
	session, ok := sm.sessions[code]
	if !ok {
		return model.Session{}, model.ErrSessionNotFound
	}

	// Проверяем, не просрочен ли токен
	if time.Now().After(session.ExpiresAt) {
		err := sm.DeleteSession(ctx, code)
		if err != nil {
			return model.Session{}, err
		}
		return model.Session{}, model.ErrSessionExpired
	}
	return session, nil
}

func (sm *SessionMemory) AddSession(ctx context.Context, session model.Session) error {
	sm.sessions[session.UUID] = session
	return nil
}

func (sm *SessionMemory) DeleteSession(ctx context.Context, code string) error {
	delete(sm.sessions, code)
	return nil
}
