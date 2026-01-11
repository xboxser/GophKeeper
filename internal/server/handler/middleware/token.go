package middleware

import (
	"context"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/service"
	"net/http"
	"time"
)

// избавляемся от проблемы коллизии ключа  в контексте
type ContextKey string

const UserIDContextKey ContextKey = "userID"

type TokenMiddleware interface {
	CheckToken(http.Handler) http.Handler
	GetUserRequest(*http.Request) model.TokenAuth
}

type tokenMiddleware struct {
	UserService    service.UserService
	TokenService   service.TokenService
	SessionService service.SessionService
}

func NewTokenMiddleware(userService service.UserService, tokenService service.TokenService, sessionService service.SessionService) *tokenMiddleware {
	return &tokenMiddleware{
		UserService:    userService,
		TokenService:   tokenService,
		SessionService: sessionService,
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

		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		// Проверяем что передали числовое значение
		uuid, err := t.TokenService.GetUser(headerValue)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		tokenAuth, err := t.SessionService.GetSession(ctx, uuid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Ищем пользователя в базе данных
		user, err := t.UserService.GetUserForID(ctx, tokenAuth.UserID)
		if user.ID < 1 || err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Добавляем userID в контекст запроса
		// На основе данного поля определяем пользователя в дальнейшем
		ctx = context.WithValue(r.Context(), UserIDContextKey, tokenAuth)
		// Передаем запрос с обновленным контекстом дальше
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserRequest - получение ID пользователя из контекста
// корректность ID проверяется в middleware
func (t *tokenMiddleware) GetUserRequest(r *http.Request) model.TokenAuth {
	tokenAuth := r.Context().Value(UserIDContextKey).(model.TokenAuth)
	return tokenAuth
}
