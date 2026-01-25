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

type CredentialHandler struct {
	CredentialService service.CredentialService
	TokenMiddleware   middleware.TokenMiddleware
}

func NewCredentialHandler(credentialService service.CredentialService, tokenMiddleware middleware.TokenMiddleware) *CredentialHandler {
	return &CredentialHandler{
		CredentialService: credentialService,
		TokenMiddleware:   tokenMiddleware,
	}
}

// Routes - поддерживает интерфейс RouteChi
func (ch *CredentialHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(ch.TokenMiddleware.CheckToken)
	r.Get("/", ch.getCredentials)
	r.Post("/", ch.addCredential)
	r.Delete("/{login}", ch.deleteCredential)
	//TODO добавить update & delete
	return r
}

// Pattern - поддерживает интерфейс RouteChi
func (_ *CredentialHandler) Pattern() string {
	return "/api/credentials"
}

func (ch *CredentialHandler) deleteCredential(w http.ResponseWriter, r *http.Request) {
	tokenAuth := ch.TokenMiddleware.GetUserRequest(r)
	if tokenAuth.UserID == 0 {
		http.Error(w, "Invalid user token", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	login := chi.URLParam(r, "login")
	if login == "" {
		http.Error(w, model.ErrCredentialEmptyLogin.Error(), http.StatusBadRequest)
		return
	}

	err := ch.CredentialService.DeleteCredential(ctx, login, tokenAuth.UserID)
	if err != nil {
		if errors.Is(err, model.ErrCredentialNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (ch *CredentialHandler) getCredentials(w http.ResponseWriter, r *http.Request) {
	tokenAuth := ch.TokenMiddleware.GetUserRequest(r)

	if tokenAuth.UserID == 0 {
		http.Error(w, "Invalid user token", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	credentials, err := ch.CredentialService.GetCredentials(ctx, tokenAuth.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(credentials)
}

func (ch *CredentialHandler) addCredential(w http.ResponseWriter, r *http.Request) {
	tokenAuth := ch.TokenMiddleware.GetUserRequest(r)

	if tokenAuth.UserID == 0 {
		http.Error(w, "Invalid user token", http.StatusUnauthorized)
		return
	}

	var credential model.CredentialAPI
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &credential); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.ValidateModelCredentialAPI(credential); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err = ch.CredentialService.AddCredential(ctx, credential, tokenAuth.UserID)
	if err != nil {
		if errors.Is(err, model.ErrCredentialDuplicate) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
