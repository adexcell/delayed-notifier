package dto

import (
	"github.com/adexcell/delayed-notifier/internal/notify/domain"
)

// GetStatusInput represents the input data required to retrieve the status of a notification.
type GetStatusInput struct {
	ID string `binding:"required" json:"id"`
}

// GetStatusOutput represents the current status of a notification.
type GetStatusOutput struct {
	Status domain.Status `json:"status"`
}
