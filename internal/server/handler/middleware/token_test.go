package middleware

import (
	"context"
	"errors"
	"gophkeeper/internal/model"
	mock_service "gophkeeper/mocks/server/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCheckToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок-сервисы
	mockUserService := mock_service.NewMockUserService(ctrl)
	mockTokenService := mock_service.NewMockTokenService(ctrl)
	mockSessionService := mock_service.NewMockSessionService(ctrl)

	middleware := NewTokenMiddleware(mockUserService, mockTokenService, mockSessionService)

	// Тест 1: Отсутствует заголовок Authorization
	t.Run("No Authorization Header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		handlerCalled := false
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
		})

		middleware.CheckToken(testHandler).ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.False(t, handlerCalled, "Обработчик не должен вызываться при отсутствии токена")
	})

	// Тест 2: Невалидный токен
	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		w := httptest.NewRecorder()

		mockTokenService.EXPECT().GetUser("Bearer invalid_token").Return("", errors.New("invalid token"))

		handlerCalled := false
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
		})

		middleware.CheckToken(testHandler).ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.False(t, handlerCalled, "Обработчик не должен вызываться при невалидном токене")
	})

	// Тест 3: Успешная проверка токена
	t.Run("Valid Token", func(t *testing.T) {
		authHeader := "Bearer valid_token"
		sessionUUID := "session-uuid"
		userID := 1

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()

		mockTokenService.EXPECT().GetUser(authHeader).Return(sessionUUID, nil)
		mockSessionService.EXPECT().GetSession(gomock.Any(), sessionUUID).Return(model.TokenAuth{UserID: userID}, nil)
		mockUserService.EXPECT().GetUserForID(gomock.Any(), userID).Return(model.User{ID: userID}, nil)

		handlerCalled := false
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenAuth := r.Context().Value(UserIDContextKey).(model.TokenAuth)
			assert.Equal(t, userID, tokenAuth.UserID)
			handlerCalled = true
		})

		middleware.CheckToken(testHandler).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, handlerCalled, "Обработчик должен вызываться при валидном токене")
	})
}

func TestGetUserRequest(t *testing.T) {
	userID := 1
	expectedTokenAuth := model.TokenAuth{UserID: userID}

	ctx := context.WithValue(context.Background(), UserIDContextKey, expectedTokenAuth)
	req := httptest.NewRequest("GET", "/", nil).WithContext(ctx)

	middleware := &tokenMiddleware{}
	result := middleware.GetUserRequest(req)

	assert.Equal(t, expectedTokenAuth, result, "Метод должен вернуть TokenAuth из контекста запроса")
	assert.Equal(t, userID, result.UserID, "UserID должен совпадать с тем, что был помещен в контекст")
}
