package http

import (
	"github.com/dezzerlol/avitotech-test-2023/internal/handlers/segment"
	"github.com/dezzerlol/avitotech-test-2023/internal/repo"
	"github.com/dezzerlol/avitotech-test-2023/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) setHTTPRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)

	segmentRepo := repo.NewSegment(s.db)

	segmentService := service.NewSegment(segmentRepo, s.worker)

	segmentHandler := segment.New(segmentService, s.logger)

	r.Post("/segment", segmentHandler.Create)
	r.Delete("/segment", segmentHandler.Delete)

	r.Post("/segment/user", segmentHandler.UpdateUserSegments)
	r.Get("/segment/user", segmentHandler.GetSegmentsForUser)
	r.Get("/segment/user/history", segmentHandler.GetUserHistory)

	r.Get("/segments/reports/{fileName}", segmentHandler.DownloadReport)

	return r
}
