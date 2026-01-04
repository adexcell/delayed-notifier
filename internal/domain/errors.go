package domain

import "errors"

var (
	// notify errors
	ErrNotFoundNotify = errors.New("not found notify")
	ErrNotifyAlreadyExisis = errors.New("notify already exists")
)
