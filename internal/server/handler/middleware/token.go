package middleware

import (
	"context"
	"gophkeeper/internal/server/service"
	"net/http"
	"time"
)

// избавляемся от проблемы коллизии ключа  в контексте
type ContextKey string

const UserIDContextKey ContextKey = "userID"

type TokenMiddleware interface {
	CheckToken(http.Handler) http.Handler
	GetUserRequest(*http.Request) int
}

type tokenMiddleware struct {
	UserService service.UserService
}

func NewTokenMiddleware(userService service.UserService) *tokenMiddleware {
	return &tokenMiddleware{
		UserService: userService,
	}
}

// CheckToken - Проверка наличия токена в запросе
func (t *tokenMiddleware) CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerValue := r.Header.Get("Authorization")

		if headerValue == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Проверяем что передали числовое значение
		userID := service.GetUserID(headerValue)
		if userID == -1 {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		// Ищем пользователя в базе данных
		user, err := t.UserService.GetUserForID(ctx, userID)
		if user.ID < 1 || err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Добавляем userID в контекст запроса
		// На основе данного поля определяем пользователя в дальнейшем
		ctx = context.WithValue(r.Context(), UserIDContextKey, user.ID)
		// Передаем запрос с обновленным контекстом дальше
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserRequest - получение ID пользователя из контекста
// корректность ID проверяется в middleware
func (t *tokenMiddleware) GetUserRequest(r *http.Request) int {
	userID := r.Context().Value(UserIDContextKey).(int)
	return userID
}
