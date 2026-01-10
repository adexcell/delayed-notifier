package rabbit

import (
	"context"
	"fmt"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/pkg/rabbit"

	"github.com/wb-go/wbf/rabbitmq"
)

const (
	exchangeName = "delayed_exchange"
	queueName    = "notifications_queue"
	routingKey   = "notification_key"
	contentType  = "application/json"
)

type NotifyQueueAdapter struct {
	client    *rabbitmq.RabbitClient
	publisher *rabbitmq.Publisher
}

func NewRabbitQueueAdapter(cfg rabbit.Config) (domain.QueueProvider, error) {
	client, err := rabbit.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init rabbitmq: %w", err)
	}

	pub := rabbitmq.NewPublisher(client, exchangeName, contentType)

	q := &NotifyQueueAdapter{
		client:    client,
		publisher: pub,
	}

	if err := q.client.DeclareQueue(
		queueName,
		exchangeName,
		routingKey,
		true,  // durable
		false, // autoDelete
		true,  // exchangeDurable
		nil,
	); err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return q, nil
}

func (q *NotifyQueueAdapter) Publish(ctx context.Context, n *domain.Notify) error {
	delay := time.Until(n.ScheduledAt)
	if delay < 0 {
		delay = 0
	}

	body, err := toRabbitDTO(n)
	if err != nil {
		return fmt.Errorf("marshal notify: %w", err)
	}

	return q.publisher.Publish(
		ctx,
		body,
		routingKey,
		rabbitmq.WithExpiration(delay),
	)
}

func (q *NotifyQueueAdapter) Close() error {
	return q.client.Close()
}
