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
	//TODO добавить update & delete
	return r
}

// Pattern - поддерживает интерфейс RouteChi
func (ch *CredentialHandler) Pattern() string {
	return "/api/credentials"
}

func (ch *CredentialHandler) getCredentials(w http.ResponseWriter, r *http.Request) {
	userID := ch.TokenMiddleware.GetUserRequest(r)

	if userID == 0 {
		http.Error(w, "Invalid user token", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	credentials, err := ch.CredentialService.GetCredentials(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(credentials)
}

func (ch *CredentialHandler) addCredential(w http.ResponseWriter, r *http.Request) {
	userID := ch.TokenMiddleware.GetUserRequest(r)

	if userID == 0 {
		http.Error(w, "Invalid user token", http.StatusUnauthorized)
		return
	}

	var credential model.Credential
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

	err = ch.CredentialService.AddCredential(ctx, credential, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
