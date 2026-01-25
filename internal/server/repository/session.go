package repository

import (
	"context"
	"gophkeeper/internal/model"
	"sync"
	"time"
)

type SessionRepository interface {
	GetSession(ctx context.Context, code string) (model.Session, error)
	AddSession(ctx context.Context, session model.Session) error
	DeleteSession(ctx context.Context, code string) error
}

type SessionMemory struct {
	sessions map[string]model.Session
	mutex    sync.Mutex
}

func NewSessionMemory() *SessionMemory {
	return &SessionMemory{
		sessions: make(map[string]model.Session),
	}
}

// GetSession - получение сессии
func (sm *SessionMemory) GetSession(ctx context.Context, code string) (model.Session, error) {
	sm.mutex.Lock()
	session, ok := sm.sessions[code]
	sm.mutex.Unlock()

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

// AddSession - добавление сессии
func (sm *SessionMemory) AddSession(_ context.Context, session model.Session) error {
	sm.mutex.Lock()
	sm.sessions[session.UUID] = session
	sm.mutex.Unlock()
	return nil
}

// DeleteSession - удаление сессии
func (sm *SessionMemory) DeleteSession(_ context.Context, code string) error {
	sm.mutex.Lock()
	delete(sm.sessions, code)
	sm.mutex.Unlock()
	return nil
}
