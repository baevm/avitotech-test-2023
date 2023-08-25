package models

import "time"

type Segment struct {
	Id        int64      `json:"id,omitempty"`
	Slug      string     `json:"slug"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
