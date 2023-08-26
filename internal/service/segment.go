package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/dezzerlol/avitotech-test-2023/internal/db/models"
	"github.com/dezzerlol/avitotech-test-2023/internal/worker"
)

type SegmentRepo interface {
	Create(ctx context.Context, segment *models.Segment) error
	DeleteBySlug(ctx context.Context, segment *models.Segment) error

	AddUserSegments(ctx context.Context, userId int64, addSegments []string, ttl int64) (int64, error)
	DeleteUserSegments(ctx context.Context, userId int64, deleteSegments []string) (int64, error)
	GetUserSegments(ctx context.Context, userId int64) ([]*models.Segment, error)

	GetUserHistory(ctx context.Context, userId int64, date time.Time) ([]*models.UserHistory, error)
}

type Segment struct {
	segmentRepo SegmentRepo
	worker      worker.TaskDistributor
}

func NewSegment(segmentRepo SegmentRepo, worker worker.TaskDistributor) *Segment {
	return &Segment{
		segmentRepo: segmentRepo,
		worker:      worker,
	}
}

func (s *Segment) Create(ctx context.Context, segment *models.Segment) error {
	return s.segmentRepo.Create(ctx, segment)
}

func (s *Segment) DeleteBySlug(ctx context.Context, segment *models.Segment) error {
	return s.segmentRepo.DeleteBySlug(ctx, segment)
}

func (s *Segment) UpdateUserSegments(
	ctx context.Context,
	userId int64,
	addSegments []string,
	ttl int64,
	deleteSegments []string,
) (
	segmentsAdded int64,
	segmentsDeleted int64,
	err error,
) {
	if len(addSegments) > 0 {
		segmentsAdded, err = s.segmentRepo.AddUserSegments(ctx, userId, addSegments, ttl)
		if err != nil {
			return segmentsAdded, segmentsDeleted, err
		}

		// Если задан TTL, то добавляем таски на удаление сегментов
		if ttl > 0 {
			for _, v := range addSegments {
				payload := worker.SegmentExpirePayload{
					UserID:      userId,
					SegmentSlug: v,
					ExpireAt:    ttl,
				}

				s.worker.ScheduleSegmentExpireTask(ctx, payload)
			}
		}
	}

	if len(deleteSegments) > 0 {
		segmentsDeleted, err = s.segmentRepo.DeleteUserSegments(ctx, userId, deleteSegments)
		if err != nil {
			return segmentsAdded, segmentsDeleted, err
		}
	}

	return segmentsAdded, segmentsDeleted, nil
}

func (s *Segment) GetUserSegments(ctx context.Context, userId int64) ([]*models.Segment, error) {
	return s.segmentRepo.GetUserSegments(ctx, userId)
}

func (s *Segment) GetUserHistory(ctx context.Context, userId, month, year int64) (string, error) {
	date := time.Date(int(year), time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	userHistory, err := s.segmentRepo.GetUserHistory(ctx, userId, date)

	if err != nil {
		return "", err
	}

	return s.generateCSVReport(userId, userHistory)
}

func (s *Segment) generateCSVReport(userId int64, userHistory []*models.UserHistory) (string, error) {
	fileName := fmt.Sprintf("./reports/%d-%d.csv", userId, time.Now().Unix())

	// Создаем файл
	csvFile, err := os.Create(fileName)

	if err != nil {
		return "", err
	}

	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)

	// Записываем слайс значения в файл
	for _, history := range userHistory {
		csvWriter.Write([]string{
			fmt.Sprintf("%d", history.UserID),
			history.SegmentSlug,
			history.Operation,
			history.ExecutedAt.Format("2006-01-02 15:04:05"),
		})
	}

	csvWriter.Flush()

	// убираем точку в начале пути
	// возвращаем ссылку в формате /reports/file_name
	return fileName[1:], nil
}
