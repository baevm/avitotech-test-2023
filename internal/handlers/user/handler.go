package user

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
}

type userService interface {
	Create(ctx context.Context) (int64, error)
}

type handler struct {
	logger  *zap.SugaredLogger
	userSvc userService
}

func NewHandler(logger *zap.SugaredLogger, userSvc userService) Handler {
	return &handler{
		logger:  logger,
		userSvc: userSvc,
	}
}
