package domain

import "errors"

// Ошибки пользователей
var (
	ErrUserNotFound           = errors.New("user not found")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrPasswordTooLong        = errors.New("password is too long")
)

// Ошибки инфраструктуры
var (
	ErrInternal            = errors.New("internal server error")
	ErrQueueFailed         = errors.New("failed to send message to queue")
	ErrCacheFailed         = errors.New("failed to process cache")
	ErrNotificationInvalid = errors.New("invalid notification data")
)
