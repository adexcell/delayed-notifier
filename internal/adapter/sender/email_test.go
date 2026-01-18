package sender

import (
	"context"
	"strings"
	"testing"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/pkg/log"
)

func TestEmailSender_Send_InvalidAddress(t *testing.T) {
	cfg := EmailConfig{
		FromEmail: "invalid-email", // Невалидный email
		FromName:  "Notifier",
	}

	sender := NewEmailSender(cfg, log.New())
	notify := &domain.Notify{
		Target: "recipient@example.com",
	}

	err := sender.Send(context.Background(), notify)

	// Ожидаем ошибку валидации адреса (go-mail проверяет RFC 5322)
	if err == nil {
		t.Fatal("expected error for invalid from address, got nil")
	}

	if !strings.Contains(err.Error(), "failed to set from address") {
		t.Errorf("expected 'failed to set from address' error, got %v", err)
	}
}

func TestEmailSender_Send_InvalidTarget(t *testing.T) {
	cfg := EmailConfig{
		FromEmail: "sender@example.com",
		FromName:  "Notifier",
	}

	sender := NewEmailSender(cfg, log.New())
	notify := &domain.Notify{
		Target: "invalid-target-email", // Невалидный targer
	}

	err := sender.Send(context.Background(), notify)

	if err == nil {
		t.Fatal("expected error for invalid target address, got nil")
	}

	if !strings.Contains(err.Error(), "failed to set to address") {
		t.Errorf("expected 'failed to set to address' error, got %v", err)
	}
}
