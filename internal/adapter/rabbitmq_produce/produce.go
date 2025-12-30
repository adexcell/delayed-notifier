package rabbitmqproduce

import (
	"context"
	"encoding/json"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type NotifyQueue struct {
	conn *amqp.Connection
}

func NewNotifyQueue(conn *amqp.Connection) *NotifyQueue {
	return &NotifyQueue{conn: conn}
}

func (q *NotifyQueue) Publish(ctx context.Context, n *domain.Notify) error {
	ch, err := q.conn.Channel()
	if err != nil {
		return domain.ErrQueueFailed
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"delayed_exchange",  // имя
		"x-delayed-message", // ТИП
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		amqp.Table{
			"x-delayed-type": "direct", // какой тип будет у "внутреннего" обмена
		},
	)
	if err != nil {
		return domain.ErrQueueFailed
	}

	delay := time.Until(n.SendAt).Milliseconds()
	if delay < 0 {
		delay = 0
	}

	body, err := json.Marshal(n)
	if err != nil {
		return domain.ErrQueueFailed
	}

	if err := ch.PublishWithContext(ctx, "delayed_exchange", "notification_key", false, false, amqp.Publishing{
		Headers:     amqp.Table{"x-delay": delay},
		ContentType: "application/json",
		Body:        body,
	}); err != nil {
		return domain.ErrQueueFailed
	}

	return nil
}

func (q *NotifyQueue) Consume(ctx context.Context) (<-chan amqp.Delivery, error) {
	ch, err := q.conn.Channel()
	if err != nil {
		return nil, domain.ErrQueueFailed
	}

	queue, err := ch.QueueDeclare(
		"notifications_queue", // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return nil, domain.ErrQueueFailed
	}

	err = ch.QueueBind(
		queue.Name,         // queue name
		"notification_key", // routing key (тот же, что в Publish)
		"delayed_exchange", // exchange
		false,
		nil,
	)

	return ch.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack (ставим false, затем вручную вызываем в воркере)
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
}
