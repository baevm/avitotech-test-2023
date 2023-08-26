package http

import (
	"github.com/dezzerlol/avitotech-test-2023/internal/handlers/segment"
	"github.com/dezzerlol/avitotech-test-2023/internal/repo"
	"github.com/dezzerlol/avitotech-test-2023/internal/service"

	_ "github.com/dezzerlol/avitotech-test-2023/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) setHTTPRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	segmentRepo := repo.NewSegment(s.db)

	segmentService := service.NewSegment(segmentRepo, s.worker)

	segmentHandler := segment.New(segmentService, s.logger)

	r.Post("/segment", segmentHandler.Create)
	r.Delete("/segment", segmentHandler.Delete)

	r.Post("/segment/user", segmentHandler.UpdateUserSegments)
	r.Get("/segment/user/{userId}", segmentHandler.GetSegmentsForUser)
	r.Get("/segment/history/{userId}", segmentHandler.GetUserHistory)

	r.Get("/segment/reports/{fileName}", segmentHandler.DownloadReport)

	return r
}
