package models

import "time"

type Segment struct {
	Slug      string     `json:"slug"`
	CreatedAt *time.Time `json:"-"`
}
