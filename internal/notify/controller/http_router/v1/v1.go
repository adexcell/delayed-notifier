package v1

import (
	"context"

	"github.com/adexcell/delayed-notifier/internal/notify/dto"
)

// Usecase defines the business logic operations for notifications.
//
//go:generate mockgen -source=handlers.go -destination=mocks/mock.go
type Usecase interface {
	CreateNotify(ctx context.Context, input dto.CreateNotifyInput) (dto.CreateNotifyOutput, error)
	GetStatus(ctx context.Context, input dto.GetStatusInput) (dto.GetStatusOutput, error)
	GetNotify(ctx context.Context, input dto.GetStatusInput) (dto.GetNotifyOutput, error)
	DeleteNotify(ctx context.Context, input dto.DeleteNotifyInput) error
}

// Handler implements HTTP handlers for notification operations.
type Handler struct {
	usecase Usecase
}

// New creates a new instance of Handler with the given usecase.
func New(usecase Usecase) *Handler {
	return &Handler{usecase: usecase}
}
