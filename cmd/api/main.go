package main

import (
	"fmt"

	"github.com/dezzerlol/avitotech-test-2023/cfg"
	"github.com/dezzerlol/avitotech-test-2023/internal/db"
	"github.com/dezzerlol/avitotech-test-2023/internal/http"
	"github.com/dezzerlol/avitotech-test-2023/internal/worker"
	"github.com/dezzerlol/avitotech-test-2023/pkg/logger"
	"github.com/hibiken/asynq"
)

// @title          Avitotech Test 2023 API
// @version         1.0
// @description     Тествое задание для стажировки Avitotech 2023

// @host      localhost:8080
// @BasePath  /
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

	redisOpts := asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%s", cfg.Get().REDIS_HOST, cfg.Get().REDIS_PORT),
	}

	distributor := worker.NewTaskDistributor(redisOpts, logger)
	processor := worker.NewTaskProcessor(redisOpts, logger, db)

	go func() {
		err := processor.Start()

		if err != nil {
			logger.Fatalf("Error starting task processor: %s", err)
		}
	}()

	server := http.New(logger, db, distributor)
	server.Run(cfg.Get().API_HOST, cfg.Get().API_PORT)
}
