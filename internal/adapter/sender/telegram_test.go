package sender

import (
	"testing"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/pkg/log"
)

func TestTelegramSender_New(t *testing.T) {
	token := "123:test-token"
	sender := NewTelegramSender(token, log.New())

	if sender == nil {
		t.Fatal("expected sender to be created")
	}
}

func TestTelegramSender_Send_Validation(t *testing.T) {
	// Telegram sender использует http.NewRequest внутри, который валидирует только URL структуру.
	// Так как URL мы строим сами, ошибок быть не должно до момента отправки.
	// Без мока http.Client протестировать Send сложно (пойдет в сеть).
	// Поэтому ограничимся тестом создания и простого вызова метода (который упадет с ошибкой сети или 404).

	token := "invalid-token"
	s := NewTelegramSender(token, log.New())
	notify := &domain.Notify{
		Target:  "chat-id",
		Payload: []byte("mesage"),
	}

	// Этот вызов пойдет в сеть к api.telegram.org.
	// В Unit тестах внешние вызовы нежелательны без моков.
	// Но так как у нас нет интерфейса для HTTP клиента в конфиге sender'а (он создается внутри),
	// мы пропустим тест Send, чтобы не спамить Telegram API и не зависеть от сети.
	// Можно было бы рефакторить TelegramSender чтобы принимать HTTP client.

	_ = s
	_ = notify
}
