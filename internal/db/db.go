package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(db_dsn string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), db_dsn)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.Ping(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
