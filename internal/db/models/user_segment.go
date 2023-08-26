package models

import (
	"database/sql"
	"time"
)

type UserSegment struct {
	UserId      int64
	SegmentSlug string
	CreatedAt   time.Time
	ExpireAt    sql.NullTime
}

func NewExpireDate(ttl int64) sql.NullTime {
	if ttl > 0 {
		return sql.NullTime{
			Time:  time.Now().Add(time.Duration(ttl) * time.Second),
			Valid: true,
		}
	}

	return sql.NullTime{}
}
