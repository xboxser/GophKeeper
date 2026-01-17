package route

import (
	"bytes"
	"context"
	"encoding/json"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/handler/middleware"
	"gophkeeper/internal/server/service"
	"gophkeeper/internal/validator"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type CardHandler struct {
	CardService     service.CardService
	TokenMiddleware middleware.TokenMiddleware
}

func NewCardHandler(cardService service.CardService, tokenMiddleware middleware.TokenMiddleware) *CardHandler {
	return &CardHandler{
		CardService:     cardService,
		TokenMiddleware: tokenMiddleware,
	}
}

// Routes - поддерживает интерфейс RouteChi
func (ch *CardHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(ch.TokenMiddleware.CheckToken)
	r.Get("/", ch.getCards)
	r.Post("/", ch.addCard)
	//TODO добавить update & delete
	return r
}

// Pattern - поддерживает интерфейс RouteChi
func (ch *CardHandler) Pattern() string {
	return "/api/card"
}

func (ch *CardHandler) getCards(w http.ResponseWriter, r *http.Request) {
	tokenAuth := ch.TokenMiddleware.GetUserRequest(r)

	if tokenAuth.UserID == 0 {
		http.Error(w, "Invalid user token", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	credentials, err := ch.CardService.GetCards(ctx, tokenAuth.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(credentials)
}

func (ch *CardHandler) addCard(w http.ResponseWriter, r *http.Request) {
	tokenAuth := ch.TokenMiddleware.GetUserRequest(r)

	if tokenAuth.UserID == 0 {
		http.Error(w, "Invalid user token", http.StatusUnauthorized)
		return
	}

	var card model.CardAPI
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &card); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.ValidateModelCardAPI(card); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err = ch.CardService.AddCard(ctx, card, tokenAuth.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
