package main

import (
	"github.com/dezzerlol/avitotech-test-2023/cfg"
	"github.com/dezzerlol/avitotech-test-2023/internal/db"
	"github.com/dezzerlol/avitotech-test-2023/internal/http"
	"github.com/dezzerlol/avitotech-test-2023/pkg/logger"
)

func main() {
	logger := logger.New()
	err := cfg.Load(".")

	if err != nil {
		logger.Fatalf("Error reading config: %s", err)
	}

	db, err := db.New(cfg.Get().DB_DSN)

	if err != nil {
		logger.Fatalf("Error starting db: %s", err)
	}

	defer db.Close()

	server := http.New(logger, db)
	server.Run(cfg.Get().HTTP_HOST, cfg.Get().HTTP_PORT)
}
