package rabbit

import (
	"context"
	"fmt"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/pkg/rabbit"
	"github.com/rabbitmq/amqp091-go"

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
		return nil, fmt.Errorf("failed to connect rabbitmq: %w", err)
	}

	pub := rabbitmq.NewPublisher(client, exchangeName, contentType)

	return &NotifyQueueAdapter{
		client:    client,
		publisher: pub,
	}, nil
}

func (q *NotifyQueueAdapter) Init() error {
	if err := q.client.DeclareQueue(
		queueName,
		exchangeName,
		routingKey,
		true,  // durable
		false, // autoDelete
		true,  // exchangeDurable
		nil,
	); err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	return nil
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

func (q *NotifyQueueAdapter) Consume(ctx context.Context, handler domain.MessageHandler) error {
	wbfHandler := func(c context.Context, d amqp091.Delivery) error {
		return handler(c, d.Body)
	}

	cfg := rabbitmq.ConsumerConfig{
		Queue:         queueName,
		ConsumerTag:   "notifier-worker",
		Workers:       5,
		PrefetchCount: 10,
	}

	consumer := rabbitmq.NewConsumer(q.client, cfg, wbfHandler)
	return consumer.Start(ctx)
}

func (q *NotifyQueueAdapter) Close() error {
	return q.client.Close()
}
