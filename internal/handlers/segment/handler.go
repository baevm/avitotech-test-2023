package segment

import (
	"net/http"

	"github.com/dezzerlol/avitotech-test-2023/internal/service"
	"go.uber.org/zap"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)

	UpdateUserSegments(w http.ResponseWriter, r *http.Request)
	GetSegmentsForUser(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	logger  *zap.SugaredLogger
	segment *service.Segment
}

func New(logger *zap.SugaredLogger, segment *service.Segment) Handler {
	return &handler{
		logger:  logger,
		segment: segment,
	}
}
