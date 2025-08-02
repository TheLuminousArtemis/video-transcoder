package main

import (
	"log/slog"
	"os"

	"github.com/hibiken/asynq"
	"github.com/theluminousartemis/video-transcoder/internal/env"
	"github.com/theluminousartemis/video-transcoder/internal/queue"
)

func main() {
	cfg := config{
		addr:       env.GetString("ADDR", ":3030"),
		asynqredis: env.GetString("ASYNQ_REDIS", "localhost:6379"),
		asynqCfg: asynqConfig{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	}

	//logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	//queuemanager
	logger.Info("connecting asynq redis client to redis", "addr", cfg.asynqredis)
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.asynqredis})
	defer client.Close()

	queueMgr := &queue.QueueManager{
		AsynqClient: client,
		AsynqServer: nil,
	}

	app := application{
		config:   cfg,
		logger:   logger,
		queueMgr: queueMgr,
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
