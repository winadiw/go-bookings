package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/winadiw/go-bookings/pkg/config"
	"github.com/winadiw/go-bookings/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	// Example of grouping for same middleware
	mux.Group(func(r chi.Router) {
		// r.Use(WriteToConsole)
		r.Get("/", handlers.Repo.Home)
	})

	mux.Get("/about", handlers.Repo.About)

	return mux

}
