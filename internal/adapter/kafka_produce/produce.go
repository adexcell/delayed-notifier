package kafkaproduce

import (
	"context"
	"encoding/json"

	"github.com/adexcell/delayed-notifier/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	l *logger.Zerolog
}

func NewProducer(writer *kafka.Writer) *Producer {
	return &Producer{writer: writer}
}

func (p *Producer) Produce(ctx context.Context, topic string, payload any) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		p.l.Error().
			Err(err).
			Str("topic", topic).
			Any("payload", payload).
			Msg("kafka producer payload marshal failed")
			
		return err
	}

	p.l.Info().
		Str("topic", topic).
		Any("payload", payload).
		Msg("kafka producer payload successfull")

	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: payloadBytes,
	})
}

func (p *Producer) Close() error {
	if err := p.writer.Close(); err != nil {
		p.l.Error().Err(err).Msg("kafka producer close failed")
		return err
	}
	return nil
}
