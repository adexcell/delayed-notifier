package domain

import "errors"

var (
	// notify errors
	ErrNotFound            = errors.New("not found notify")
	ErrNotifyAlreadyExists = errors.New("notify already exists")
)
