package repo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dezzerlol/avitotech-test-2023/internal/db/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Segment struct {
	DB *pgxpool.Pool
}

func NewSegment(db *pgxpool.Pool) *Segment {
	return &Segment{DB: db}
}

var (
	ErrSegmentNotFound = errors.New("segment not found")
)

func (r Segment) Create(ctx context.Context, segment *models.Segment) error {
	query := `
		INSERT INTO segments (slug)
		VALUES ($1)
		RETURNING created_at
	`

	args := []any{
		segment.Slug,
	}

	err := r.DB.
		QueryRow(ctx, query, args...).
		Scan(&segment.CreatedAt)

	return err
}

func (r Segment) DeleteBySlug(ctx context.Context, segment *models.Segment) error {
	query := `
		DELETE FROM segments
		WHERE slug = $1
	`

	args := []any{&segment.Slug}

	ct, err := r.DB.Exec(ctx, query, args...)

	if ct.RowsAffected() == 0 {
		return ErrSegmentNotFound
	}

	return err
}

func (r Segment) CreateUser(ctx context.Context) (int64, error) {
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

func (r Segment) CheckUserExist(ctx context.Context, userId int64) (bool, error) {
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

func (r Segment) GetUserSegments(ctx context.Context, userId int64) ([]*models.Segment, error) {
	query := `
		SELECT slug
		FROM segments s
		JOIN user_segments us
		on s.slug = us.segment_slug
		WHERE us.user_id = $1
	`

	args := []any{userId}

	rows, err := r.DB.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	var segments []*models.Segment

	for rows.Next() {
		var segment models.Segment

		err := rows.Scan(
			&segment.Slug,
		)

		if err != nil {
			return nil, err
		}

		segments = append(segments, &segment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return segments, nil
}

func (r Segment) AddUserSegments(ctx context.Context, userId int64, addSegments []string, ttl int64) (int64, error) {
	var sb strings.Builder

	// Сначала получаем slug сегментов, которые нужно добавить
	// Потом вставляем в user_segments
	sb.WriteString(`
	INSERT INTO user_segments (segment_slug, user_id, expire_at)
	SELECT s.slug, $1, $2
	FROM segments s
	WHERE s.slug IN (
	`)

	args := []any{userId, models.NewExpireDate(ttl)}

	// Готовим аргументы для запроса
	// Добавялем 3 потому что первый аргумент это userId
	for i, slug := range addSegments {
		args = append(args, slug)
		sb.WriteString(fmt.Sprintf("$%d", i+3))

		if i != len(addSegments)-1 {
			sb.WriteString(",")
		}
	}

	// Закрываем строку запроса
	// Добавляем пропуск конфликта, если сегмент уже есть у пользователя
	sb.WriteString(") ON CONFLICT DO NOTHING")

	query := sb.String()

	ct, err := r.DB.Exec(ctx, query, args...)

	return ct.RowsAffected(), err
}

func (r Segment) AddRndUsersSegment(ctx context.Context, slug string, percent int8) error {
	// Ищем процент рандомных пользователей
	// И создаем записи в user_segments
	// В случае если запись уже существует пропускаем
	query := `
	INSERT INTO user_segments (segment_slug, user_id)
	SELECT s.slug, u.id 
	FROM users u
	JOIN segments s ON s.slug = $1
	ORDER BY random() 
	LIMIT (SELECT count(1) FROM users) * ($2/100.0)
	ON CONFLICT DO NOTHING`

	args := []any{slug, percent}

	_, err := r.DB.Exec(ctx, query, args...)

	return err
}

func (r Segment) DeleteUserSegments(ctx context.Context, userId int64, deleteSegments []string) (int64, error) {
	var sb strings.Builder

	// Сначала получаем slug сегментов, которые нужно удалить
	// Потом удаляем из user_segments
	sb.WriteString(`
	DELETE FROM user_segments us
	WHERE us.user_id = $1
	AND us.segment_slug IN (
		SELECT s.slug
		FROM segments s
		WHERE s.slug IN (
	`)

	args := []any{userId}

	// Готовим аргументы для запроса
	// Добавялем 2 потому что первый аргумент это userId
	for i, slug := range deleteSegments {
		args = append(args, slug)
		sb.WriteString(fmt.Sprintf("$%d", i+2))

		if i != len(deleteSegments)-1 {
			sb.WriteString(",")
		}
	}

	// Закрываем строку запроса
	sb.WriteString("))")

	query := sb.String()

	ct, err := r.DB.Exec(ctx, query, args...)

	return ct.RowsAffected(), err
}

func (r Segment) GetUserHistory(ctx context.Context, userId int64, date time.Time) ([]*models.UserHistory, error) {
	query := `
		SELECT user_id, segment_slug, operation, executed_at
		FROM user_segment_history
		WHERE user_id = $1
		AND date_part('year', executed_at) = $2 
		AND date_part('month', executed_at) = $3`

	args := []any{
		userId,
		date.Year(),
		date.Month(),
	}

	rows, err := r.DB.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	var history []*models.UserHistory

	for rows.Next() {
		var userHistory models.UserHistory

		err := rows.Scan(
			&userHistory.UserID,
			&userHistory.SegmentSlug,
			&userHistory.Operation,
			&userHistory.ExecutedAt,
		)

		if err != nil {
			return nil, err
		}

		history = append(history, &userHistory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}
