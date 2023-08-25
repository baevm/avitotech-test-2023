package models

import "time"

type UserSegment struct {
	UserId      int64
	SegmentSlug string
	CreatedAt   time.Time
}
