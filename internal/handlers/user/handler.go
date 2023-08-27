package user

import (
	"net/http"

	"github.com/dezzerlol/avitotech-test-2023/internal/service"
	"go.uber.org/zap"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	logger  *zap.SugaredLogger
	userSvc *service.User
}

func New(logger *zap.SugaredLogger, userSvc *service.User) Handler {
	return &handler{
		logger:  logger,
		userSvc: userSvc,
	}
}
