package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeName string = "delayed_exchange"
	queueName    string = "notifications_queue"
	routingKey   string = "notification_key"
)

type NotifyQueueAdapter struct {
	conn *amqp.Connection
}

func NewRabbitQueueAdapter(conn *amqp.Connection) *NotifyQueueAdapter {
	return &NotifyQueueAdapter{conn: conn}
}

func (q *NotifyQueueAdapter) Init(ctx context.Context) error {
	ch, err := q.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open connection channel rabbitmq: %w", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName,        // name
		"x-delayed-message", // exchange type
		true,                // durable
		false,               // auto delete
		false,               // internal
		false,               // no wait
		amqp.Table{
			"x-delayed-type": "direct",
		},
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange rabbitmq: %w", err)
	}

	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = ch.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)

	return nil
}

func (q *NotifyQueueAdapter) Publish(ctx context.Context, n *domain.Notify) error {
	ch, err := q.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open connection channel rabbitmq: %w", err)
	}
	defer ch.Close()

	delay := max(time.Until(n.ScheduledAt).Milliseconds(), 0)

	dto := toRabbitDTO(n)

	body, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("marshal notify: %w", err)
	}

	if err := ch.PublishWithContext(ctx, exchangeName, routingKey, false, false, amqp.Publishing{
		Headers:     amqp.Table{"x-delay": delay},
		ContentType: "application/json",
		Body:        body,
	}); err != nil {
		return fmt.Errorf("failed to publish rabbitmq: %w", err)
	}

	return nil
}

func (q *NotifyQueueAdapter) Consume(ctx context.Context, handler domain.MessageHandler) error {
	ch, err := q.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open connection channel rabbitmq: %w", err)
	}

	outputChan, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed channel consume: %w", err)
	}

	return nil
}

func (q *NotifyQueueAdapter) Close(ch *amqp.Channel) {
	ch.Close()
}
