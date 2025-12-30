package schedule

import (
	"context"

	kafkaproduce "github.com/adexcell/delayed-notifier/internal/adapter/kafka_produce"
	rabbitmqproduce "github.com/adexcell/delayed-notifier/internal/adapter/rabbitmq_produce"

	"github.com/adexcell/delayed-notifier/internal/domain"
)

type RabbitMq interface {
	Publish(ctx context.Context, n *domain.Notify) error
}

type Kafka interface {
	Produce(ctx context.Context, topic string, payload any) error
}

type Usecase struct {
	rabbitMq RabbitMq
	kafka    Kafka
}

func New(rabbitMq *rabbitmqproduce.NotifyQueue, kafka *kafkaproduce.Producer)

func (c *Usecase) Schedule(ctx context.Context, n *domain.Notify) error {

	if err := c.rabbitMq.Publish(ctx, n); err != nil {
		return err
	}

	return c.kafka.Produce(ctx, domain.Topic.Notify, n)
}
