package domain

import (
	"fmt"
	"time"

	"github.com/adexcell/delayed-notifier/pkg/validate"
	"github.com/google/uuid"
)

type Channel string

const (
	ChannelEmail      Channel = "email"
	ChannelTelegram   Channel = "telegram"
	defaultMaxRetries         = 3
)

type Notify struct {
	ID          uuid.UUID
	Channel     Channel   `validate:"required"`
	Recipient   string    `validate:"required"`
	Subject     *string   // optionally for email, may be nil
	Body        string    `validate:"required"`
	ScheduledAt time.Time `validate:"required"`
	Status      Status
	RetryCount  int
	MaxRetries  int
	LastError   *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewNotify(
	channel Channel,
	recipient string,
	subject *string,
	body string,
	scheduledAt time.Time,
) (Notify, error) {
	n := Notify{
		ID:          uuid.New(),
		Channel:     channel,
		Recipient:   recipient,
		Subject:     subject,
		Body:        body,
		ScheduledAt: scheduledAt,
		Status:      Pending,
		RetryCount:  0,
		MaxRetries:  defaultMaxRetries,
		CreatedAt:   time.Now().UTC(),
	}

	if err := validate.Validate(n); err != nil {
		return Notify{}, fmt.Errorf("n.Validate: %w", err)
	}

	return n, nil
}
