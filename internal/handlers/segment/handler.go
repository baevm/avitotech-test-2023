package segment

import (
	"context"
	"net/http"

	"github.com/dezzerlol/avitotech-test-2023/cfg"
	"github.com/dezzerlol/avitotech-test-2023/internal/db/models"
	"go.uber.org/zap"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)

	UpdateUserSegments(w http.ResponseWriter, r *http.Request)
	GetSegmentsForUser(w http.ResponseWriter, r *http.Request)

	GetUserHistory(w http.ResponseWriter, r *http.Request)
	DownloadReport(w http.ResponseWriter, r *http.Request)
}

//go:generate mockgen -destination=mocks/mock_segment.go -package=mocks github.com/dezzerlol/avitotech-test-2023/internal/handlers/segment SegmentService
type SegmentService interface {
	Create(ctx context.Context, segment *models.Segment) error
	DeleteBySlug(ctx context.Context, segment *models.Segment) error
	GetUserSegments(ctx context.Context, userId int64) ([]*models.Segment, error)
	GetUserHistory(ctx context.Context, userId, month, year int64) (string, error)
	UpdateUserSegments(ctx context.Context, userId int64, addSegments []string, ttl int64, deleteSegments []string) (segmentsAdded int64, segmentsDeleted int64, err error)
}

type handler struct {
	logger     *zap.SugaredLogger
	segmentSvc SegmentService
	config     *cfg.Config
}

func NewHandler(logger *zap.SugaredLogger, segment SegmentService) Handler {
	return &handler{
		logger:     logger,
		segmentSvc: segment,
		config:     cfg.Get(),
	}
}
