package user

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
}

//go:generate mockgen -destination=mocks/mock_user.go -package=mocks github.com/dezzerlol/avitotech-test-2023/internal/handlers/user UserService
type UserService interface {
	Create(ctx context.Context) (int64, error)
}

type handler struct {
	logger  *zap.SugaredLogger
	userSvc UserService
}

func NewHandler(logger *zap.SugaredLogger, userSvc UserService) Handler {
	return &handler{
		logger:  logger,
		userSvc: userSvc,
	}
}
