package repo

import "errors"

var (
	UniqueConstraintViolation = "23505"
)

var (
	// Common errors
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")

	// Segment errors
	ErrSegmentAlreadyExists = errors.New("segment already exists")
	ErrSegmentNotFound      = errors.New("segment not found")

	// User errors
	ErrUserNotFound = errors.New("user not found")
)
