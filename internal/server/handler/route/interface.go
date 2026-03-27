package route

import "github.com/go-chi/chi/v5"

// RouteChi - interface для элементов роутинга chi
type RouteChi interface {
	// Routes - возвращает роутер chi.
	Routes() chi.Router
	// Pattern - возвращает паттерн роутера chi.
	Pattern() string
}
