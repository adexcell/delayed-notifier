// Package mailer provides implementations for sending notifications via email
package mailer

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/rs/zerolog/log"
)

const defaultSendTimeout = 10 * time.Second

var (
	// ErrInvalidHeader mean contains invalid header characters.
	ErrInvalidHeader = errors.New("contains invalid header characters")
)

// Config holds the configuration parameters for the EmailSender.
type Config struct {
	Host        string        `envconfig:"MAILER_HOST"`
	Port        int           `envconfig:"MAILER_PORT"`
	Email       string        `envconfig:"MAILER_EMAIL"`
	Password    string        `envconfig:"MAILER_PASSWORD"`
	UseTLS      bool          `envconfig:"MAILER_USE_TLS"`
	SendTimeout time.Duration `envconfig:"MAILER_SEND_TIMEOUT"`
}

// EmailSender handles sending notifications via SMTP.
type EmailSender struct {
	cfg  Config
	auth smtp.Auth
}

// NewEmailSender creates a new instance of EmailSender.
func NewEmailSender(cfg Config) *EmailSender {
	auth := smtp.PlainAuth("", cfg.Email, cfg.Password, cfg.Host)

	return &EmailSender{
		cfg:  cfg,
		auth: auth,
	}
}

// Send sends a notification email to the recipient specified in the domain model.
func (s *EmailSender) Send(ctx context.Context, notify domain.Notify) (err error) {
	from, to, msg, err := s.prepareMessage(notify)
	if err != nil {
		return err
	}

	client, closeClient, err := s.connect(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := closeClient(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	if err := s.sendMessage(client, from, to, msg); err != nil {
		return err
	}

	log.Debug().Msg("EmailSender: successfully sent")

	return nil
}

func (s *EmailSender) connect(ctx context.Context) (*smtp.Client, func() error, error) {
	deadline := sendDeadline(ctx, s.cfg.SendTimeout)

	timeout := time.Until(deadline)
	if timeout <= 0 {
		return nil, nil, context.DeadlineExceeded
	}

	addr := net.JoinHostPort(s.cfg.Host, fmt.Sprintf("%d", s.cfg.Port))
	tlsConfig := s.tlsConfig()

	netDialer := &net.Dialer{
		Timeout: timeout,
	}

	var conn net.Conn
	var err error

	if s.cfg.Port == 465 || (s.cfg.UseTLS && s.cfg.Port != 587) {
		tlsDialer := &tls.Dialer{
			NetDialer: netDialer,
			Config:    tlsConfig,
		}

		conn, err = tlsDialer.DialContext(ctx, "tcp", addr)
	} else {
		conn, err = netDialer.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("dial smtp: %w", err)
	}

	if err := conn.SetDeadline(deadline); err != nil {
		_ = conn.Close()
		return nil, nil, fmt.Errorf("set smtp deadline: %w", err)
	}

	client, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		_ = conn.Close()
		return nil, nil, fmt.Errorf("smtp client: %w", err)
	}

	closeClient := func() error {
		quitErr := client.Quit()
		closeErr := conn.Close()

		return errors.Join(
			wrapError("smtp quit", quitErr),
			wrapError("close smtp connection", closeErr),
		)
	}

	return client, closeClient, nil
}

func (s *EmailSender) sendMessage(client *smtp.Client, from, to *mail.Address, msg []byte) error {
	if s.cfg.UseTLS && s.cfg.Port != 465 {
		if err := client.StartTLS(s.tlsConfig()); err != nil {
			return fmt.Errorf("starttls: %w", err)
		}
	}

	if err := client.Auth(s.auth); err != nil {
		return fmt.Errorf("[mail-sender] auth: %w", err)
	}

	if err := client.Mail(from.Address); err != nil {
		return fmt.Errorf("[mail-sender] mail from: %w", err)
	}

	if err := client.Rcpt(to.Address); err != nil {
		return fmt.Errorf("[mail-sender] rcpt to %s: %w", to.Address, err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("data stream: %w", err)
	}

	if err := writeSMTPData(w, msg); err != nil {
		return err
	}

	return nil
}

func (s *EmailSender) tlsConfig() *tls.Config {
	return &tls.Config{
		ServerName: s.cfg.Host,
		MinVersion: tls.VersionTLS12,
	}
}

func wrapError(message string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", message, err)
}

func (s *EmailSender) prepareMessage(notify domain.Notify) (*mail.Address, *mail.Address, []byte, error) {
	from, err := parseEmailAddress(s.cfg.Email)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parse sender email: %w", err)
	}

	to, err := parseEmailAddress(notify.RecipientEmail)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parse recipient email: %w", err)
	}
	if hasHeaderInjection(notify.Subject) {
		return nil, nil, nil, ErrInvalidHeader
	}

	msg, err := buildMessage(from.String(), to.String(), notify.Subject, notify.Body)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("build message: %w", err)
	}

	return from, to, msg, nil
}

func buildMessage(from, to, subject, body string) ([]byte, error) {
	var msg bytes.Buffer

	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", mime.QEncoding.Encode("UTF-8", subject)))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	msg.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")

	qp := quotedprintable.NewWriter(&msg)
	if _, err := qp.Write([]byte(body)); err != nil {
		_ = qp.Close()
		return nil, fmt.Errorf("encode body: %w", err)
	}

	if err := qp.Close(); err != nil {
		return nil, fmt.Errorf("close body encoder: %w", err)
	}

	return msg.Bytes(), nil
}

func writeSMTPData(w io.WriteCloser, msg []byte) error {
	writeErr := error(nil)

	if _, err := w.Write(msg); err != nil {
		writeErr = fmt.Errorf("write data: %w", err)
	}

	closeErr := w.Close()
	if closeErr != nil {
		closeErr = fmt.Errorf("close data stream: %w", closeErr)
	}

	return errors.Join(writeErr, closeErr)
}

func parseEmailAddress(value string) (*mail.Address, error) {
	if hasHeaderInjection(value) {
		return nil, ErrInvalidHeader
	}

	return mail.ParseAddress(value)
}

func hasHeaderInjection(value string) bool {
	return strings.ContainsAny(value, "\r\n")
}

func sendDeadline(ctx context.Context, timeout time.Duration) time.Time {
	if deadline, ok := ctx.Deadline(); ok {
		return deadline
	}

	if timeout <= 0 {
		timeout = defaultSendTimeout
	}

	return time.Now().Add(timeout)
}
