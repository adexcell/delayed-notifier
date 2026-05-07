package domain

import "errors"

var (
	// ErrNotFound is returned when a notification is not found in the repository.
	ErrNotFound = errors.New("notify not found")
	// ErrInvalidID is returned when the provided notification ID is invalid.
	ErrInvalidID = errors.New("invalid id")
	// ErrScheduledTime is returned when the scheduled time exceeds the maximum allowed delay (approx. 50 days).
	ErrScheduledTime = errors.New("scheduled time too far in the future")
)
