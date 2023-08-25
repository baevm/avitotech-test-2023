package service

import (
	"context"

	"github.com/dezzerlol/avitotech-test-2023/internal/db/models"
)

type SegmentRepo interface {
	Create(ctx context.Context, segment *models.Segment) error
	DeleteBySlug(ctx context.Context, segment *models.Segment) error

	AddUserSegments(ctx context.Context, userId int64, addSegments []string) (int64, error)
	DeleteUserSegments(ctx context.Context, userId int64, deleteSegments []string) (int64, error)
	GetUserSegments(ctx context.Context, userId int64) ([]*models.Segment, error)
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
