package v1

import (
	"context"

	"github.com/adexcell/delayed-notifier/internal/notify/dto"
)

type Usecase interface {
	CreateNotify(ctx context.Context, input dto.CreateNotifyInput) (dto.CreateNotifyOutput, error)
	GetNotify(ctx context.Context, input dto.GetNotifyInput) (dto.GetNotifyOutput, error)
	DeleteNotify(ctx context.Context, input dto.DeleteNotifyInput) error
	UpdateNotify(ctx context.Context, input dto.UpdateNotifyInput) error
}

type Handler struct {
	usecase Usecase
}

func New(usecase Usecase) *Handler {
	return &Handler{usecase: usecase}
}
