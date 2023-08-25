package user

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Handler interface {
}

type handler struct {
	logger *zap.SugaredLogger
}

func New(logger *zap.SugaredLogger, db *pgxpool.Pool) Handler {
	return &handler{
		logger: logger,
	}
}
