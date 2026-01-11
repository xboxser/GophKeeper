package service

import (
	"context"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/repository"
	"time"

	"github.com/google/uuid"
)

type SessionService interface {
	GetSession(ctx context.Context, code string) (model.TokenAuth, error)
	AddSession(ctx context.Context, tokenAuth model.TokenAuth) (model.Session, error)
}

type sessionService struct {
	SessionRepository repository.SessionRepository
	ExpireTime        time.Duration
}

func NewSessionService(sessionRepository repository.SessionRepository, expireTime time.Duration) *sessionService {
	return &sessionService{
		SessionRepository: sessionRepository,
		ExpireTime:        expireTime,
	}
}

func (s *sessionService) GetSession(ctx context.Context, code string) (model.TokenAuth, error) {
	session, err := s.SessionRepository.GetSession(ctx, code)

	if err != nil {
		return model.TokenAuth{}, err
	}

	return model.TokenAuth{UserID: session.UserID}, nil
}

func (s *sessionService) AddSession(ctx context.Context, tokenAuth model.TokenAuth) (model.Session, error) {
	uuid := uuid.New().String()
	if uuid == "" {
		return model.Session{}, model.ErrSessionGenerateUUID
	}

	session := model.Session{
		UUID:      uuid,
		UserID:    tokenAuth.UserID,
		ExpiresAt: time.Now().Add(s.ExpireTime),
	}

	err := s.SessionRepository.AddSession(ctx, session)
	if err != nil {
		return model.Session{}, err
	}
	return session, nil
}
