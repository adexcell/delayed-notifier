package domain

import (
	"fmt"
	"math"
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/domain/validation"
	"github.com/google/uuid"
)

const (
	defaultMaxRetries = 3
)

// Notify represents a notification entity in the system.
type Notify struct {
	ID             uuid.UUID
	RecipientEmail string    `validate:"required,email"`
	Subject        string    `validate:"required,lte=255"`
	Body           string    `validate:"required"`
	ScheduledAt    time.Time `validate:"required,gt"`
	Status         Status
	RetryCount     int
	MaxRetries     int
	LastError      *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewNotify creates a new Notify entity with the given parameters and validates it.
func NewNotify(
	email string,
	subject string,
	body string,
	scheduledAt time.Time,
) (Notify, error) {
	n := Notify{
		ID:             uuid.New(),
		RecipientEmail: email,
		Subject:        subject,
		Body:           body,
		ScheduledAt:    scheduledAt.UTC(),
		Status:         StatusPending,
		RetryCount:     0,
		MaxRetries:     defaultMaxRetries,
		CreatedAt:      time.Now().UTC(),
	}

	delay := n.ScheduledAt.UnixMilli() - time.Now().UnixMilli()
	if delay > math.MaxUint32 {
		return Notify{}, ErrScheduledTime
	}

	err := validation.Validate(n)
	if err != nil {
		return Notify{}, fmt.Errorf("n.Validate: %w", err)
	}

	return n, nil
}
