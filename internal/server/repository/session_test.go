package repository

import (
	"context"
	"gophkeeper/internal/model"
	"testing"
	"time"
)

func TestGetSession(t *testing.T) {
	ctx := context.Background()
	repo := NewSessionMemory()

	activeSession := model.Session{
		UUID:      "active-session",
		UserID:    1,
		ExpiresAt: time.Now().Add(1 * time.Hour), // Активная сессия
	}

	expiredSession := model.Session{
		UUID:      "expired-session",
		UserID:    2,
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Просроченная сессия
	}

	// Добавляем сессии в хранилище
	_ = repo.AddSession(ctx, activeSession)
	_ = repo.AddSession(ctx, expiredSession)

	t.Run("get existing active session", func(t *testing.T) {
		session, err := repo.GetSession(ctx, "active-session")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if session.UserID != 1 {
			t.Errorf("Expected UserID 1, got %d", session.UserID)
		}
	})

	t.Run("get non-existing session", func(t *testing.T) {
		_, err := repo.GetSession(ctx, "non-existing")
		if err != model.ErrSessionNotFound {
			t.Errorf("Expected ErrSessionNotFound, got %v", err)
		}
	})

	t.Run("get expired session", func(t *testing.T) {
		_, err := repo.GetSession(ctx, "expired-session")
		if err != model.ErrSessionExpired {
			t.Errorf("Expected ErrSessionExpired, got %v", err)
		}

		// Проверяем, что просроченная сессия была удалена
		_, err = repo.GetSession(ctx, "expired-session")
		if err != model.ErrSessionNotFound {
			t.Errorf("Expected ErrSessionNotFound after deletion, got %v", err)
		}
	})
}
