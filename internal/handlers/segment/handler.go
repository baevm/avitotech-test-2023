package segment

import (
	"net/http"

	"github.com/dezzerlol/avitotech-test-2023/cfg"
	"github.com/dezzerlol/avitotech-test-2023/internal/service"
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

type handler struct {
	logger  *zap.SugaredLogger
	segment *service.Segment
	config  *cfg.Config
}

func New(segment *service.Segment, logger *zap.SugaredLogger) Handler {
	return &handler{
		logger:  logger,
		segment: segment,
		config:  cfg.Get(),
	}
}
