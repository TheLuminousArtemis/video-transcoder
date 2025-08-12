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
	config := &config{
		redisAddr: env.GetString("REDIS_ADDR", "localhost:6379"),
		redisqueueCfg: asynqConfig{
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

	//redis queue
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.redisAddr},
		asynq.Config{
			Concurrency: config.redisqueueCfg.Concurrency,
			Queues:      config.redisqueueCfg.Queues,
		},
	)
	logger.Info("successfully connected to redis")

	//redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.redisCfg.addr,
		Password: config.redisCfg.password,
		DB:       config.redisCfg.db,
	})

	//queuemanager
	queueMgr := queue.QueueManager{
		AsynqClient: nil,
		AsynqServer: server,
	}

	app := application{
		logger:   logger,
		config:   config,
		queueMgr: queueMgr,
		rdb:      rdb,
	}

	err := runAsynqWorker(&app)
	if err != nil {
		logger.Error(err.Error())
	}

}
