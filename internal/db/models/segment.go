package models

import "time"

type Segment struct {
	Slug        string     `json:"slug"`
	UserPercent int8       `json:"user_percent,omitempty"`
	CreatedAt   *time.Time `json:"-"`
}
