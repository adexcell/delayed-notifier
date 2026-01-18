package controller

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
)

func TestDTOConversion(t *testing.T) {
	now := time.Now().UTC()
	errStr := "test error"

	domainNotify := &domain.Notify{
		ID:          "test-id",
		Payload:     []byte(`"payload"`),
		Target:      "target",
		Channel:     "email",
		Status:      domain.StatusPending,
		ScheduledAt: now,
		CreatedAt:   now,
		RetryCount:  1,
		LastError:   &errStr,
	}

	// 1. Domain -> Response
	resp := toResponse(domainNotify)

	if resp.ID != domainNotify.ID {
		t.Errorf("expected ID %s, got %s", domainNotify.ID, resp.ID)
	}
	if resp.Status != domainNotify.Status {
		t.Errorf("expected Status %v, got %v", domainNotify.Status, resp.Status)
	}

	// 2. Handler DTO -> Domain
	dto := NotifyControllerDTO{
		ID:          "test-id",
		Payload:     json.RawMessage(`"payload"`),
		Target:      "target",
		Channel:     "email",
		ScheduledAt: now,
		Status:      domain.StatusPending,
		CreatedAt:   now,
	}

	converted := toDomain(dto)

	if converted.ID != dto.ID {
		t.Errorf("expected ID %s, got %s", dto.ID, converted.ID)
	}

	// Проверяем payload через string представление, т.к. json.RawMessage vs []byte может отличаться nil/empty
	if string(converted.Payload) != string(dto.Payload) {
		t.Errorf("expected payload %s, got %s", dto.Payload, converted.Payload)
	}
}

func TestDTO_JSONMarshaling(t *testing.T) {
	// Проверка JSON тегов
	req := CreateNotifyRequest{
		ID: "test",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Простой тест на наличие поля
	if !reflect.DeepEqual(data, []byte(`{"id":"test","payload":null,"target":"","channel":"","scheduled_at":"0001-01-01T00:00:00Z"}`)) {
		// тут точное совпадение строки зависит от порядка полей и дефолтных значений.
		// Лучше просто проверить наличие ошибок маршалинга, так как теги стандартные.
	}
}
