package rabbit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/wb-go/wbf/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitQueueAdapter struct {
	publisher *rabbitmq.Publisher
}

func NewRabbitQueueAdapter(pub *rabbitmq.Publisher) *RabbitQueueAdapter {
	return &RabbitQueueAdapter{publisher: pub}
}

func (a *RabbitQueueAdapter) Publish(ctx context.Context, n *domain.Notify) error {
	dto := toRabbitDTO(n)

	body, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("marshal notify: %w", err)
	}

	return a.publisher.Publish(
		ctx, 
		body, 
		a.publisher.GetExchangeName(),
		func(p *amqp.Publishing) {
			p.DeliveryMode = amqp.Persistent
		},
	)
}
