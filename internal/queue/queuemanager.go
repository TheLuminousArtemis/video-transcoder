package queue

import "github.com/hibiken/asynq"

type QueueManager struct {
	AsynqClient *asynq.Client
	AsynqServer *asynq.Server
}
