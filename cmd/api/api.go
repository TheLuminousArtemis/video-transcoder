package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type config struct {
	addr string
}

type application struct {
	config config
	logger *slog.Logger
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/health", app.health)
	r.Post("/upload", app.uploadVideo)
	return r
}

func (app *application) start(router http.Handler) error {

	srv := http.Server{
		Addr:    app.config.addr,
		Handler: router,
	}
	return srv.ListenAndServe()
}
