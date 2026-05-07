package usecase

import (
	"context"
	"testing"

	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDeleteNotify_Success(t *testing.T) {
	pgMock := &MockPostgres{
		DeleteNotifyFunc: func(ctx context.Context, notifyID uuid.UUID) error {
			return nil
		},
	}

	uc := New(pgMock, &MockRedis{}, &MockRedisWriter{}, &MockRabbitWriter{})

	input := dto.DeleteNotifyInput{ID: uuid.New().String()}
	err := uc.DeleteNotify(context.Background(), input)

	require.NoError(t, err)
}

func TestDeleteNotify_InvalidID(t *testing.T) {
	uc := New(&MockPostgres{}, &MockRedis{}, &MockRedisWriter{}, &MockRabbitWriter{})

	input := dto.DeleteNotifyInput{ID: "not-a-uuid"}
	err := uc.DeleteNotify(context.Background(), input)

	require.Error(t, err)
}
