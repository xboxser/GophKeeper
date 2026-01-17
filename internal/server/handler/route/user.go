package route

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/handler/middleware"
	"gophkeeper/internal/server/service"
	"gophkeeper/internal/validator"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// UserHandler - обработчик запросов связанных с пользователем
// Поддерживает интерфейс RouteChi
type UserHandler struct {
	UserService     service.UserService
	TokenMiddleware middleware.TokenMiddleware
}

func NewUserHandler(userService service.UserService, tokenMiddleware middleware.TokenMiddleware) *UserHandler {
	return &UserHandler{
		UserService:     userService,
		TokenMiddleware: tokenMiddleware,
	}
}

func (u *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/register", u.register)
	r.Post("/login", u.login)

	// Группа маршрутов использующих токен у user
	protected := r.With(u.TokenMiddleware.CheckToken)
	protected.Get("/code", u.code)
	return r
}
func (u *UserHandler) Pattern() string {
	return "/api/user"
}

func (u *UserHandler) register(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user model.APIUser

	if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.ValidateModelRegistrationUserAPI(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	token, err := u.UserService.Register(ctx, user)
	if err != nil {
		if errors.Is(err, model.ErrLoginBusy) {
			http.Error(w, "This login is busy", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)

}

func (u *UserHandler) login(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user model.APIUser
	if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.ValidateModelUserAPI(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	token, err := u.UserService.Login(ctx, user)
	if err != nil {
		if errors.Is(err, model.ErrLoginIncorrect) {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)

}

// code - Возвращаем хеш мастер кода
func (u *UserHandler) code(w http.ResponseWriter, r *http.Request) {
	tokenAuth := u.TokenMiddleware.GetUserRequest(r)

	if tokenAuth.UserID == 0 {
		http.Error(w, "Invalid user token", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	user, err := u.UserService.GetUserForID(ctx, tokenAuth.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	code := model.Code{Hash: user.Code}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(code)

}
