package route

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/service"
	"gophkeeper/internal/validator"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// UserHandler - обработчик запросов связанных с пользователем
// Поддерживает интерфейс RouteChi
type UserHandler struct {
	UserService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (u *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/register", u.register)
	return r
}
func (u *UserHandler) Pattern() string {
	return "/api/user"
}

func (u *UserHandler) register(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

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

	userID, err := u.UserService.Register(ctx, user)
	if err != nil {
		if errors.Is(err, model.ErrLoginBusy) {
			http.Error(w, "This login is busy", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := service.BuildJWTString(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)

}
