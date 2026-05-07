package usecase

import (
	"context"
	"fmt"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/adexcell/delayed-notifier/internal/notify/metrics"
	"github.com/adexcell/delayed-notifier/pkg/otel/tracer"
)

// CreateNotify creates a new notification and schedules its delivery via RabbitMQ.
func (u *NotifyUsecase) CreateNotify(ctx context.Context, input dto.CreateNotifyInput) (dto.CreateNotifyOutput, error) {
	ctx, span := tracer.Start(ctx, "notify usecase CreateNotify")
	defer span.End()

	var output dto.CreateNotifyOutput

	notify, err := domain.NewNotify(input.RecipientEmail, input.Subject, input.Body, input.ScheduledAt)
	if err != nil {
		return output, fmt.Errorf("domain.NewNotify: %w", err)
	}

	err = u.postgres.CreateNotify(ctx, notify)
	if err != nil {
		return output, fmt.Errorf("u.postgres.CreateNotify: %w", err)
	}

	delivery := dto.ToDelivery(notify)
	u.asyncRabbitWriter.Send(ctx, delivery)

	output.ID = notify.ID
	metrics.NotificationsCreatedTotal.Inc()

	return output, nil
}
