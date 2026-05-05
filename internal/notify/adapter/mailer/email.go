// Package mailer provides implementations for sending notifications via email
package mailer

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"mime/quotedprintable"
	"net"
	"net/smtp"
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
)

type Config struct {
	Host        string        `envconfig:"MAILER_HOST"`
	Port        int           `envconfig:"MAILER_PORT"`
	Email       string        `envconfig:"MAILER_EMAIL"`
	Password    string        `envconfig:"MAILER_PASSWORD"`
	UseTLS      bool          `envconfig:"MAILER_USE_TLS"`
	SendTimeout time.Duration `envconfig:"MAILER_SEND_TIMEOUT"`
}

type EmailSender struct {
	cfg  Config
	auth smtp.Auth
}

func NewEmailSender(cfg Config) *EmailSender {
	auth := smtp.PlainAuth("", cfg.Email, cfg.Password, cfg.Host)
	return &EmailSender{
		cfg:  cfg,
		auth: auth,
	}
}

func (s *EmailSender) Send(ctx context.Context, notify domain.Notify) error {
	subject := notify.Subject
	to := notify.RecipientEmail

	// 1. Формируем MIME-сообщение с поддержкой UTF-8 и Quoted-Printable
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("From: %s\r\n", s.cfg.Email))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", to))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")

	// Кодируем тело письма для безопасности передачи через SMTP
	qp := quotedprintable.NewWriter(&buf)
	if _, err := qp.Write([]byte(notify.Body)); err != nil {
		return fmt.Errorf("encode body: %w", err)
	}
	qp.Close()
	msg := buf.String()

	// 2. Управление таймаутами через context
	deadline, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		deadline = time.Now().Add(10 * time.Second) // fallback
	}
	timeout := time.Until(deadline)
	if timeout <= 0 {
		return context.DeadlineExceeded
	}

	// 3. Подключение к серверу
	addr := fmt.Sprintf("%s.%d", s.cfg.Host, s.cfg.Port)
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return fmt.Errorf("dial smtp: %w", err)
	}
	defer conn.Close()
	conn.SetDeadline(deadline)

	client, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Quit()

	// 4. STARTTLS (если требуется)
	if s.cfg.UseTLS {
		if err := client.StartTLS(&tls.Config{ServerName: s.cfg.Host}); err != nil {
			return fmt.Errorf("starttls: %w", err)
		}
	}

	// 5. Аутентификация и отправка
	if err := client.Auth(s.auth); err != nil {
		return fmt.Errorf("[mail-sender] auth: %w", err)
	}
	if err := client.Mail(s.cfg.Email); err != nil {
		return fmt.Errorf("[mail-sender] mail from: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("[mail-sender] rcpt to %s: %w", addr, err)

	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("data stream: %w", err)
	}
	if _, err := w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("write data: %w", err)
	}
	return w.Close()
}
