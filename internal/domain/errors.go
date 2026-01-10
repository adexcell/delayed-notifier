package domain

import "errors"

var (
	// notify errors
	ErrNotFound            = errors.New("not found notify")
	ErrNotifyAlreadyExisis = errors.New("notify already exists")
)
