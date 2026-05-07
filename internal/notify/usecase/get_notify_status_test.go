package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetNotify_FromRedis(t *testing.T) {
	redisMock := &MockRedis{
		GetNotifyStatusFunc: func(ctx context.Context, key string) (domain.Status, error) {
			return domain.StatusProcessing, nil
		},
	}
	pgMock := &MockPostgres{} // Should not be called

	uc := New(pgMock, redisMock, &MockRedisWriter{}, &MockRabbitWriter{})

	input := dto.GetStatusInput{ID: uuid.New().String()}
	output, err := uc.GetStatus(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, domain.StatusProcessing, output.Status)
}

func TestGetNotify_FromPostgres(t *testing.T) {
	redisMock := &MockRedis{
		GetNotifyStatusFunc: func(ctx context.Context, key string) (domain.Status, error) {
			return domain.Status(0), errors.New("cache miss")
		},
	}
	pgMock := &MockPostgres{
		GetNotifyStatusByIDFunc: func(ctx context.Context, notifyID uuid.UUID) (domain.Status, error) {
			return domain.StatusSent, nil
		},
	}

	redisWriterMock := &MockRedisWriter{}
	var savedTask domain.NotifyStatusTask
	redisWriterMock.SendFunc = func(ctx context.Context, task domain.NotifyStatusTask) {
		savedTask = task
	}

	uc := New(pgMock, redisMock, redisWriterMock, &MockRabbitWriter{})

	id := uuid.New().String()
	input := dto.GetStatusInput{ID: id}
	output, err := uc.GetStatus(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, domain.StatusSent, output.Status)

	// Verify it queued the result to Redis
	require.Equal(t, id, savedTask.Key)
	require.Equal(t, domain.StatusSent.String(), savedTask.Value)
}
