package handler

import (
	"gophkeeper/internal/server/handler/route"

	"github.com/go-chi/chi/v5"
)

type ChiHandler struct {
	Router            *chi.Mux
	credentialHandler *route.CredentialHandler
}

func NewChiHandler() *ChiHandler {
	return &ChiHandler{
		Router: chi.NewRouter(),
	}
}

func (ch *ChiHandler) AddRoutes(route route.RouteChi) {
	ch.Router.Mount(route.Pattern(), route.Routes())
}
