package main

import (
	"log/slog"
	"os"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/theluminousartemis/video-transcoder/internal/env"
	"github.com/theluminousartemis/video-transcoder/internal/queue"
)

func main() {
	cfg := config{
		addr:      env.GetString("ADDR", ":3030"),
		redisAddr: env.GetString("ASYNQ_REDIS", "localhost:6379"),
		asynqCfg: asynqConfig{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
		redisCfg: redisConfig{
			addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
			password: env.GetString("REDIS_PASSWORD", ""),
			db:       env.GetInt("REDIS_DB", 0),
		},
	}

	//logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	//queuemanager
	logger.Info("connecting asynq redis client to redis", "addr", cfg.redisAddr)
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.redisAddr})
	defer client.Close()

	queueMgr := &queue.QueueManager{
		AsynqClient: client,
		AsynqServer: nil,
	}

	//redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.redisCfg.addr,
		Password: cfg.redisCfg.password,
		DB:       cfg.redisCfg.db,
	})

	app := application{
		config:   cfg,
		logger:   logger,
		queueMgr: queueMgr,
		rdb:      rdb,
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
