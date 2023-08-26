package worker

import (
	"context"

	"github.com/dezzerlol/avitotech-test-2023/internal/repo"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type TaskProcessor interface {
	Start() error
	ProcessSegmentExpireTask(ctx context.Context, task *asynq.Task) error
}

type SegmentRepo interface {
	DeleteUserSegments(ctx context.Context, userId int64, deleteSegments []string) (int64, error)
}

type RedisTaskProcessor struct {
	server      *asynq.Server
	logger      *zap.SugaredLogger
	segmentRepo SegmentRepo
}

func NewTaskProcessor(r asynq.RedisClientOpt, logger *zap.SugaredLogger, db *pgxpool.Pool) TaskProcessor {
	server := asynq.NewServer(r, asynq.Config{
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			logger.Errorw(
				"error handling task",
				"task_type: ", task.Type(),
				"err: ", err,
			)
		}),
		Logger: logger,
	})

	segmentRepo := repo.NewSegment(db)

	return &RedisTaskProcessor{
		server:      server,
		logger:      logger,
		segmentRepo: segmentRepo,
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(SegmentExpireTaskType, p.ProcessSegmentExpireTask)

	return p.server.Run(mux)
}
