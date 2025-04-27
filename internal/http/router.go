package http

import (
	"net/http"

	"github.com/dezzerlol/avitotech-test-2023/internal/handlers/segment"
	"github.com/dezzerlol/avitotech-test-2023/internal/handlers/user"
	"github.com/dezzerlol/avitotech-test-2023/internal/repo"
	"github.com/dezzerlol/avitotech-test-2023/internal/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/dezzerlol/avitotech-test-2023/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	reqTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "segment_app_request_total",
			Help: "Total number of requests received",
		},
	)
)

func totalRequestsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqTotal.Inc()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) setHTTPRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)
	r.Use(totalRequestsMiddleware)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	r.Handle("/metrics", promhttp.Handler())

	segmentRepo := repo.NewSegmentRepo(s.db)
	userRepo := repo.NewUserRepo(s.db)

	segmentService := service.NewSegmentSvc(s.worker, segmentRepo, userRepo)
	userService := service.NewUserSvc(userRepo)

	segmentHandler := segment.NewHandler(s.logger, segmentService)
	userHandler := user.NewHandler(s.logger, userService)

	// Создание пользователя
	r.Post("/user", userHandler.Create)

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
