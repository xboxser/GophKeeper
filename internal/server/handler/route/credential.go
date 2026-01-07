package route

import (
	"encoding/json"
	"gophkeeper/internal/server/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CredentialHandler struct {
	CredentialService service.CredentialService
}

func NewCredentialHandler(credentialService service.CredentialService) *CredentialHandler {
	return &CredentialHandler{
		CredentialService: credentialService,
	}
}

// Routes - поддерживает интерфейс RouteChi
func (ch *CredentialHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", ch.getCredentials)
	return r
}

// Pattern - поддерживает интерфейс RouteChi
func (ch *CredentialHandler) Pattern() string {
	return "/api/credentials"
}

func (ch *CredentialHandler) getCredentials(w http.ResponseWriter, r *http.Request) {
	credentials, err := ch.CredentialService.GetCredentials()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(credentials)
}
