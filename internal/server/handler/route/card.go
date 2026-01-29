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
	r.Delete("/{id}", ch.deleteCard)
	r.Put("/", ch.updateCard)
	return r
}

// Pattern - поддерживает интерфейс RouteChi
func (_ *CardHandler) Pattern() string {
	return "/api/card"
}

func (ch *CardHandler) deleteCard(w http.ResponseWriter, r *http.Request) {
	tokenAuth := ch.TokenMiddleware.GetUserRequest(r)
	if tokenAuth.UserID == 0 {
		// TODO вынести проверку в отдельный метод и заменить в других файлах
		http.Error(w, "Invalid user token", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, model.ErrCardEmptyID.Error(), http.StatusBadRequest)
		return
	}

	err := ch.CardService.DeleteCard(ctx, id, tokenAuth.UserID)
	if err != nil {
		if errors.Is(err, model.ErrCardNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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

	err = ch.CardService.AddCard(ctx, &card, tokenAuth.UserID)
	if err != nil {
		if errors.Is(err, model.ErrCardDuplicate) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ch *CardHandler) updateCard(w http.ResponseWriter, r *http.Request) {
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

	err = ch.CardService.UpdateCard(ctx, &card, tokenAuth.UserID)
	if err != nil {
		if errors.Is(err, model.ErrCardNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
