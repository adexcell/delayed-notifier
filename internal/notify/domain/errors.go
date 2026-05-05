package domain

import "errors"

var (
	ErrNotFound  = errors.New("notify not found")
	ErrInvalidID = errors.New("invalid id")
	ErrScheduledTime = errors.New("The planned time must be earlier than 50 days")
)
