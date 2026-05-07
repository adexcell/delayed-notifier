package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/adexcell/delayed-notifier/pkg/otel"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type MockPostgres struct {
	CreateNotifyFunc        func(ctx context.Context, notify domain.Notify) error
	GetNotifyStatusByIDFunc func(ctx context.Context, notifyID uuid.UUID) (domain.Status, error)
	GetNotifyByIDFunc       func(ctx context.Context, notifyID uuid.UUID) (domain.Notify, error)
	DeleteNotifyFunc        func(ctx context.Context, notifyID uuid.UUID) error
}

func (m *MockPostgres) CreateNotify(ctx context.Context, notify domain.Notify) error {
	if m.CreateNotifyFunc != nil {
		return m.CreateNotifyFunc(ctx, notify)
	}
	return nil
}

func (m *MockPostgres) GetNotifyStatusByID(ctx context.Context, notifyID uuid.UUID) (domain.Status, error) {
	if m.GetNotifyStatusByIDFunc != nil {
		return m.GetNotifyStatusByIDFunc(ctx, notifyID)
	}
	return domain.Status(0), nil
}

func (m *MockPostgres) GetNotifyByID(ctx context.Context, notifyID uuid.UUID) (domain.Notify, error) {
	if m.GetNotifyByIDFunc != nil {
		return m.GetNotifyByIDFunc(ctx, notifyID)
	}
	return domain.Notify{}, nil
}

func (m *MockPostgres) DeleteNotify(ctx context.Context, notifyID uuid.UUID) error {
	if m.DeleteNotifyFunc != nil {
		return m.DeleteNotifyFunc(ctx, notifyID)
	}
	return nil
}

type MockRabbitWriter struct {
	SendFunc func(ctx context.Context, delivery dto.Delivery)
}

func (m *MockRabbitWriter) Send(ctx context.Context, delivery dto.Delivery) {
	if m.SendFunc != nil {
		m.SendFunc(ctx, delivery)
	}
}

type MockRedis struct {
	GetNotifyStatusFunc func(ctx context.Context, key string) (domain.Status, error)
}

func (m *MockRedis) GetNotifyStatus(ctx context.Context, key string) (domain.Status, error) {
	if m.GetNotifyStatusFunc != nil {
		return m.GetNotifyStatusFunc(ctx, key)
	}
	return domain.Status(0), errors.New("not found")
}

type MockRedisWriter struct {
	SendFunc func(ctx context.Context, task domain.NotifyStatusTask)
}

func (m *MockRedisWriter) Send(ctx context.Context, task domain.NotifyStatusTask) {
	if m.SendFunc != nil {
		m.SendFunc(ctx, task)
	}
}

func TestCreateNotify_Success(t *testing.T) {
	otel.SilentModeInit()
	pg := &MockPostgres{}
	rabbit := &MockRabbitWriter{}

	var createdNotify domain.Notify
	pg.CreateNotifyFunc = func(ctx context.Context, notify domain.Notify) error {
		createdNotify = notify
		return nil
	}

	var sentDelivery dto.Delivery
	rabbit.SendFunc = func(ctx context.Context, delivery dto.Delivery) {
		sentDelivery = delivery
	}

	uc := New(pg, &MockRedis{}, &MockRedisWriter{}, rabbit)

	input := dto.CreateNotifyInput{
		RecipientEmail: "test@example.com",
		Subject:        "Test Subject",
		Body:           "Test Body",
		ScheduledAt:    time.Now().Add(24 * time.Hour),
	}

	output, err := uc.CreateNotify(context.Background(), input)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, output.ID)

	// Verify Postgres was called
	require.Equal(t, output.ID, createdNotify.ID)
	require.Equal(t, input.RecipientEmail, createdNotify.RecipientEmail)

	// Verify RabbitMQ was called
	require.Equal(t, output.ID, sentDelivery.ID)
	require.Equal(t, input.RecipientEmail, sentDelivery.RecipientEmail)
}

func TestCreateNotify_InvalidInput(t *testing.T) {
	otel.SilentModeInit()
	uc := New(&MockPostgres{}, &MockRedis{}, &MockRedisWriter{}, &MockRabbitWriter{})

	input := dto.CreateNotifyInput{
		RecipientEmail: "invalid-email", // invalid email
		Subject:        "Test Subject",
		Body:           "Test Body",
		ScheduledAt:    time.Now().Add(24 * time.Hour),
	}

	_, err := uc.CreateNotify(context.Background(), input)
	require.Error(t, err) // domain.NewNotify should fail validation
}
