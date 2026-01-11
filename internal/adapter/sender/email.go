package sender

import (
	"context"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/pkg/log"
)

type EmailSender struct {
	log log.Log
}

func NewEmailSender(log log.Log) domain.Sender {
	return &EmailSender{log: log}
}

func (s *EmailSender) Send(ctx context.Context, n *domain.Notify) error {
	s.log.Info().
        Str("target", n.Target).
        Str("payload", string(n.Payload)).
        Msg("[EMAIL] Sending message")
	return nil
}
