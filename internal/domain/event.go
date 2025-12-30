package domain

import "context"

type EventSender interface {
	Send(ctx context.Context, topic string, payload any) error
}
