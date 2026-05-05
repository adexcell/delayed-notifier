package v1

import (
	"context"

	"github.com/adexcell/delayed-notifier/internal/notify/dto"
)

//go:generate mockgen -source=handlers.go -destination=mocks/mock.go
type Usecase interface {
	CreateNotify(ctx context.Context, input dto.CreateNotifyInput) (dto.CreateNotifyOutput, error)
	GetNotify(ctx context.Context, input dto.GetNotifyInput) (dto.GetNotifyOutput, error)
	DeleteNotify(ctx context.Context, input dto.DeleteNotifyInput) error
}

type Handler struct {
	usecase Usecase
}

func New(usecase Usecase) *Handler {
	return &Handler{usecase: usecase}
}
