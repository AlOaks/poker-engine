package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
)

func Router(db *sqlx.DB) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	return router
}
