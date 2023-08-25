package models

import "time"

type Segment struct {
	Id        int64      `json:"-"`
	Slug      string     `json:"slug"`
	CreatedAt *time.Time `json:"-"`
}
