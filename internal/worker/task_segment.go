package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dezzerlol/avitotech-test-2023/internal/db/models"
	"github.com/hibiken/asynq"
)

type SegmentExpirePayload struct {
	UserID      int64
	SegmentSlug string
	ExpireAt    int64 // seconds
}

const (
	SegmentExpireTaskType = "segment:expire"
)

func (d *RedisTaskDistributor) ScheduleSegmentExpireTask(ctx context.Context, payload SegmentExpirePayload, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	task := asynq.NewTask(SegmentExpireTaskType, jsonPayload, opts...)

	processAt := models.NewExpireDate(payload.ExpireAt)

	info, err := d.client.EnqueueContext(ctx, task, asynq.ProcessAt(processAt.Time))

	d.logger.Infow(
		"task enqueue",
		"task_type", info.Type,
		"task_id", info.ID,
		"queue", info.Queue,
	)

	return err
}

func (p *RedisTaskProcessor) ProcessSegmentExpireTask(ctx context.Context, task *asynq.Task) error {
	var payload SegmentExpirePayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	_, err := p.segmentRepo.DeleteUserSegments(ctx, payload.UserID, []string{payload.SegmentSlug})

	if err != nil {
		return fmt.Errorf("segmentService.UpdateUserSegments failed: %v: %w", err, asynq.SkipRetry)
	}

	p.logger.Infow(
		"task processed",
		"task_type", task.Type(),
	)

	return nil
}
