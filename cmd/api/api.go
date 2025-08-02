package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/theluminousartemis/video-transcoder/internal/queue"
)

type config struct {
	addr       string
	asynqredis string
	asynqCfg   asynqConfig
}

type application struct {
	config   config
	logger   *slog.Logger
	queueMgr *queue.QueueManager
}

type asynqConfig struct {
	Concurrency int
	Queues      map[string]int
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", app.health)
		r.Post("/upload", app.uploadVideo)
	})
	return r
}

func (app *application) start(router http.Handler) error {

	srv := http.Server{
		Addr:    app.config.addr,
		Handler: router,
	}
	return srv.ListenAndServe()
}
