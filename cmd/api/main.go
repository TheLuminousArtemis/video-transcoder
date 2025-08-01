package main

import (
	"log/slog"
	"os"

	"github.com/theluminousartemis/video-transcoder/internal/env"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":3030"),
	}

	//logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	app := application{
		config: cfg,
		logger: logger,
	}
	//mux
	mux := app.mount()

	//starting the server
	logger.Info("starting the server", "addr", app.config.addr)
	err := app.start(mux)
	if err != nil {
		logger.Error("error starting server", "addr", app.config.addr)
	}

}
