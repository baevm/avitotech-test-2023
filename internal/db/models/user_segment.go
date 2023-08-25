package models

import "time"

type UserSegment struct {
	UserId    int64
	SegmentId int64
	CreatedAt time.Time
}
