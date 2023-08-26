package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/dezzerlol/avitotech-test-2023/internal/db/models"
)

type SegmentRepo interface {
	Create(ctx context.Context, segment *models.Segment) error
	DeleteBySlug(ctx context.Context, segment *models.Segment) error

	AddUserSegments(ctx context.Context, userId int64, addSegments []string) (int64, error)
	DeleteUserSegments(ctx context.Context, userId int64, deleteSegments []string) (int64, error)
	GetUserSegments(ctx context.Context, userId int64) ([]*models.Segment, error)

	GetUserHistory(ctx context.Context, userId int64, date time.Time) ([]*models.UserHistory, error)
}

type Segment struct {
	segmentRepo SegmentRepo
}

func NewSegment(segmentRepo SegmentRepo) *Segment {
	return &Segment{
		segmentRepo: segmentRepo,
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
	deleteSegments []string,
) (
	segmentsAdded int64,
	segmentsDeleted int64,
	err error) {

	if len(addSegments) > 0 {
		segmentsAdded, err = s.segmentRepo.AddUserSegments(ctx, userId, addSegments)
		if err != nil {
			return segmentsAdded, segmentsDeleted, err
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

func (s *Segment) GetUserHistory(ctx context.Context, userId int64, date time.Time) (string, error) {
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
