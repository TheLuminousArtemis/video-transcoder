package main

import (
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/theluminousartemis/video-transcoder/internal/queue"
)

type config struct {
	redisAddr     string
	redisqueueCfg asynqConfig
}

type asynqConfig struct {
	Concurrency int
	Queues      map[string]int
}

type application struct {
	logger   *slog.Logger
	config   *config
	queueMgr queue.QueueManager
}

func runAsynqWorker(app *application) error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(queue.TypeVideoTranscode, app.HandleTranscodeTask)
	if err := app.queueMgr.AsynqServer.Run(mux); err != nil {
		return err
	}
	return nil
}
