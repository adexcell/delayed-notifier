package sender

import (
	"context"
	"fmt"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/pkg/log"
	"github.com/wneessen/go-mail"
)

type EmailConfig struct {
	SMTPHost     string `mapstructure:"smtp_host"`
	SMTPPort     int    `mapstructure:"smtp_port"`
	SMTPUsername string `mapstructure:"smtp_username"`
	SMTPPassword string `mapstructure:"smtp_password"`
	FromEmail    string `mapstructure:"from_email"`
	FromName     string `mapstructure:"from_name"`
}

type EmailSender struct {
	config EmailConfig
	log    log.Log
}

func NewEmailSender(config EmailConfig, log log.Log) domain.Sender {
	return &EmailSender{
		config: config,
		log:    log,
	}
}

func (s *EmailSender) Send(ctx context.Context, n *domain.Notify) error {
	m := mail.NewMsg()

	if err := m.From(fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)); err != nil {
		return fmt.Errorf("failed to set from address: %w", err)
	}

	if err := m.To(n.Target); err != nil {
		return fmt.Errorf("failed to set to address: %w", err)
	}

	m.Subject("Delayed Notification")
	m.SetBodyString(mail.TypeTextPlain, string(n.Payload))

	client, err := mail.NewClient(
		s.config.SMTPHost,
		mail.WithPort(s.config.SMTPPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(s.config.SMTPUsername),
		mail.WithPassword(s.config.SMTPPassword),
		mail.WithTLSPolicy(mail.TLSMandatory),
	)
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}

	sendCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := client.DialAndSendWithContext(sendCtx, m); err != nil {
		s.log.Error().
			Err(err).
			Str("target", n.Target).
			Msg("failed to send email")
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.log.Info().
		Str("target", n.Target).
		Msg("email sent successfully")
	return nil
}
