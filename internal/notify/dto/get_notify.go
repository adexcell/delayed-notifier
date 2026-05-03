package dto

import (
	"github.com/adexcell/delayed-notifier/internal/notify/domain"
)

type GetNotifyInput struct {
	ID string `binding:"required" json:"id"`
}

type GetNotifyOutput struct {
	Status domain.Status `json:"status"`
}
