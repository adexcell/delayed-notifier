package postgres

import (
	"log"
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
		Payload:     []byte("payload"),
		Target:      "target",
		Channel:     "email",
		Status:      domain.StatusPending,
		ScheduledAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
		RetryCount:  1,
		LastError:   &errStr,
	}

	// 1. Domain -> DTO
	dto := toPostgresDTO(domainNotify)

	if dto.ID != domainNotify.ID {
		t.Errorf("expected ID %s, got %s", domainNotify.ID, dto.ID)
	}

	// 2. DTO -> Domain
	convertedNotify := toDomain(dto)

	if !reflect.DeepEqual(convertedNotify, domainNotify) {
		// Log fields for debugging
		log.Printf("Expected: %+v", domainNotify)
		log.Printf("Got:      %+v", convertedNotify)
		t.Error("converted notify does not match original notify")
	}
}
