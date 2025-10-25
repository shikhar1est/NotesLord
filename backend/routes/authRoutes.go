package routes

import (
	"github.com/go-chi/chi/v5"
	"noteslord/handlers"
	"noteslord/database"
)

// RegisterAuthRoutes handles /register and /login endpoints
func RegisterAuthRoutes(r chi.Router) {
	r.Post("/register", handlers.Register(database.DB))
	r.Post("/login", handlers.Login(database.DB))
}
