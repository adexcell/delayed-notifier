package domain

import (
	"context"
	"time"
)

type Notify struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Message   string    `json:"message" binding:"required"`
	Status    string    `json:"status"`
	SendAt    time.Time `json:"send_at" binding:"required"`
	CreatedAt time.Time `json:"created_at" binding:"required"`
}

type NotificationRepository interface {
	Create(ctx context.Context, n *Notify) error
}
