package usecase

import (
	"context"

	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/adexcell/delayed-notifier/pkg/otel/tracer"
)

func (u *NotifyUsecase) CreateNotify(ctx context.Context, input dto.CreateNotifyInput) (dto.CreateNotifyOutput, error) {
	ctx, span := tracer.Start(ctx, "notify usecase CreateNotify")
	defer span.End()

	var output dto.CreateNotifyOutput

	// notify, err := domain.NewNotify(input.Channel, input.Recipient, input.Subject, input.Body, input.ScheduledAt)
	// if err != nil {
	// 	return output, fmt.Errorf("domain.NewNotify: %w", err)
	// }

	// output.ID = notify.ID

	return output, nil
}
