package dto

import (
	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/google/uuid"
)

type GetNotifyInput struct {
	ID uuid.UUID `json:"id"`
}

type GetNotifyOutput struct {
	Status domain.Status `json:"status"`
}
