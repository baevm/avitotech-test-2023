package models

import "time"

type UserHistory struct {
	ID          int64     `json:"-"`
	SegmentSlug string    `json:"segment_slug"`
	UserID      int64     `json:"user_id"`
	Operation   string    `json:"operation"`
	ExecutedAt  time.Time `json:"executed_at"`
}
