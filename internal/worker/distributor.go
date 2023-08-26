package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type TaskDistributor interface {
	ScheduleSegmentExpireTask(ctx context.Context, payload SegmentExpirePayload, opts ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
	logger *zap.SugaredLogger
}

func NewTaskDistributor(redis asynq.RedisClientOpt, logger *zap.SugaredLogger) TaskDistributor {
	client := asynq.NewClient(redis)

	return &RedisTaskDistributor{
		client: client,
		logger: logger,
	}
}
