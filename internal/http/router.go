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

	// Создание сегмента
	r.Post("/segment", segmentHandler.Create)
	// Удаление сегмента
	r.Delete("/segment", segmentHandler.Delete)

	// Добавление/удаление сегментов у пользователя
	r.Post("/segment/user", segmentHandler.UpdateUserSegments)
	// Получение всех сегментов пользователя
	r.Get("/segment/user/{userId}", segmentHandler.GetSegmentsForUser)
	// Получение ссылки на отчет по сегментам пользователя
	r.Get("/segment/history/{userId}", segmentHandler.GetUserHistory)
	// Скачивание отчета пользователя по сегментам
	r.Get("/segment/reports/{fileName}", segmentHandler.DownloadReport)

	return r
}
