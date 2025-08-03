package main

import (
	"log/slog"
	"os"

	"github.com/hibiken/asynq"
	"github.com/theluminousartemis/video-transcoder/internal/env"
	"github.com/theluminousartemis/video-transcoder/internal/queue"
)

func main() {
	config := &config{
		redisAddr: env.GetString("ASYNQ_REDIS", "localhost:6379"),
		redisqueueCfg: asynqConfig{
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

	//redis queue
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.redisAddr},
		asynq.Config{
			Concurrency: config.redisqueueCfg.Concurrency,
			Queues:      config.redisqueueCfg.Queues,
		},
	)

	//queuemanager
	queueMgr := queue.QueueManager{
		AsynqClient: nil,
		AsynqServer: server,
	}

	app := application{
		logger:   logger,
		config:   config,
		queueMgr: queueMgr,
	}

	err := runAsynqWorker(&app)
	if err != nil {
		logger.Error(err.Error())
	}

}
