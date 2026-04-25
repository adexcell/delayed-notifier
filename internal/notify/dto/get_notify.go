package dto

import (
	"github.com/adexcell/delayed-notifier/internal/notify/domain"
)

type GetNotifyInput struct {
	ID string `json:"id"`
}

type GetNotifyOutput struct {
	Status domain.Status `json:"status"`
}
