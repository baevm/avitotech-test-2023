package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	DB *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *User {
	return &User{DB: db}
}

func (r User) CreateUser(ctx context.Context) (int64, error) {
	query := `
		INSERT INTO users DEFAULT VALUES
		RETURNING id
	`

	var userId int64

	err := r.DB.
		QueryRow(ctx, query).
		Scan(&userId)

	return userId, err
}

func (r User) CheckUserExist(ctx context.Context, userId int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE id = $1
		)
	`

	args := []any{userId}

	var exists bool

	err := r.DB.
		QueryRow(ctx, query, args...).
		Scan(&exists)

	return exists, err
}
