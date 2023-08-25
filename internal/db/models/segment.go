package models

import "time"

type Segment struct {
	Id        int64
	Slug      string
	CreatedAt time.Time
}
