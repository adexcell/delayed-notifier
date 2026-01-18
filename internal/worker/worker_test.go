package worker

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/adexcell/delayed-notifier/config"
	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/internal/mocks"
	"github.com/adexcell/delayed-notifier/pkg/log"
	"go.uber.org/mock/gomock"
)

func TestNotifyConsumer_Handle_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks.NewMockNotifyPostgres(ctrl)
	mockRedis := mocks.NewMockNotifyRedis(ctrl)
	mockQueue := mocks.NewMockQueueProvider(ctrl)
	mockSender := mocks.NewMockSender(ctrl)

	cfg := config.NotifierConfig{
		MaxRetries: 3,
	}

	senders := map[string]domain.Sender{
		"email": mockSender,
	}

	consumer := NewNotifyConsumer(cfg, mockPostgres, mockQueue, mockRedis, senders, log.New())

	ctx := context.Background()
	notify := &domain.Notify{
		ID:         "test-id-123",
		Target:     "test@example.com",
		Channel:    "email",
		Payload:    []byte("Test message"),
		Status:     domain.StatusPending,
		RetryCount: 0,
	}

	payload, _ := json.Marshal(NotifyWorkerDTO{
		ID:         notify.ID,
		Target:     notify.Target,
		Channel:    notify.Channel,
		Payload:    notify.Payload,
		RetryCount: 0,
	})

	// Expect: получение из Redis (кеш промах)
	mockRedis.EXPECT().
		Get(ctx, notify.ID).
		Return(nil, errors.New("not in cache")).
		Times(1)

	// Expect: получение из БД
	mockPostgres.EXPECT().
		GetNotifyByID(ctx, notify.ID).
		Return(notify, nil).
		Times(1)

	// Expect: успешная отправка
	mockSender.EXPECT().
		Send(ctx, gomock.Any()).
		Return(nil).
		Times(1)

	// Expect: обновление статуса на Sent
	mockPostgres.EXPECT().
		UpdateStatus(ctx, notify.ID, domain.StatusSent, nil, 0, nil).
		Return(nil).
		Times(1)

	// Act
	err := consumer.Handle(ctx, payload)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestNotifyConsumer_Handle_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks.NewMockNotifyPostgres(ctrl)
	mockRedis := mocks.NewMockNotifyRedis(ctrl)
	mockQueue := mocks.NewMockQueueProvider(ctrl)

	cfg := config.NotifierConfig{MaxRetries: 3}
	senders := map[string]domain.Sender{}

	consumer := NewNotifyConsumer(cfg, mockPostgres, mockQueue, mockRedis, senders, log.New())

	ctx := context.Background()
	invalidPayload := []byte("invalid json")

	// Act
	err := consumer.Handle(ctx, invalidPayload)

	// Assert - должно вернуть nil (не ретраить невалидный JSON)
	if err != nil {
		t.Errorf("expected nil error for invalid JSON, got %v", err)
	}
}

func TestNotifyConsumer_Handle_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks.NewMockNotifyPostgres(ctrl)
	mockRedis := mocks.NewMockNotifyRedis(ctrl)
	mockQueue := mocks.NewMockQueueProvider(ctrl)

	cfg := config.NotifierConfig{MaxRetries: 3}
	senders := map[string]domain.Sender{}

	consumer := NewNotifyConsumer(cfg, mockPostgres, mockQueue, mockRedis, senders, log.New())

	ctx := context.Background()
	notifyID := "non-existent-id"

	payload, _ := json.Marshal(NotifyWorkerDTO{
		ID:      notifyID,
		Target:  "test@example.com",
		Channel: "email",
	})

	// Expect: кеш промах
	mockRedis.EXPECT().
		Get(ctx, notifyID).
		Return(nil, errors.New("not in cache")).
		Times(1)

	// Expect: не найдено в БД
	mockPostgres.EXPECT().
		GetNotifyByID(ctx, notifyID).
		Return(nil, domain.ErrNotFound).
		Times(1)

	// Act
	err := consumer.Handle(ctx, payload)

	// Assert - должно вернуть nil (не ретраить несуществующее уведомление)
	if err != nil {
		t.Errorf("expected nil error for not found, got %v", err)
	}
}

func TestNotifyConsumer_Handle_AlreadySent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks.NewMockNotifyPostgres(ctrl)
	mockRedis := mocks.NewMockNotifyRedis(ctrl)
	mockQueue := mocks.NewMockQueueProvider(ctrl)

	cfg := config.NotifierConfig{MaxRetries: 3}
	senders := map[string]domain.Sender{}

	consumer := NewNotifyConsumer(cfg, mockPostgres, mockQueue, mockRedis, senders, log.New())

	ctx := context.Background()
	notify := &domain.Notify{
		ID:      "test-id-123",
		Status:  domain.StatusSent, // Уже отправлено
		Target:  "test@example.com",
		Channel: "email",
	}

	payload, _ := json.Marshal(NotifyWorkerDTO{
		ID:      notify.ID,
		Target:  notify.Target,
		Channel: notify.Channel,
	})

	// Expect: получение из Redis
	mockRedis.EXPECT().
		Get(ctx, notify.ID).
		Return(notify, nil).
		Times(1)

	// Act
	err := consumer.Handle(ctx, payload)

	// Assert - должно пропустить уже отправленное
	if err != nil {
		t.Errorf("expected nil error for already sent, got %v", err)
	}
}

func TestNotifyConsumer_Handle_SendFailure_Retry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks.NewMockNotifyPostgres(ctrl)
	mockRedis := mocks.NewMockNotifyRedis(ctrl)
	mockQueue := mocks.NewMockQueueProvider(ctrl)
	mockSender := mocks.NewMockSender(ctrl)

	cfg := config.NotifierConfig{
		MaxRetries: 3,
	}

	senders := map[string]domain.Sender{
		"email": mockSender,
	}

	consumer := NewNotifyConsumer(cfg, mockPostgres, mockQueue, mockRedis, senders, log.New())

	ctx := context.Background()
	notify := &domain.Notify{
		ID:         "test-id-123",
		Target:     "test@example.com",
		Channel:    "email",
		Payload:    []byte("Test message"),
		Status:     domain.StatusPending,
		RetryCount: 0,
	}

	payload, _ := json.Marshal(NotifyWorkerDTO{
		ID:         notify.ID,
		Target:     notify.Target,
		Channel:    notify.Channel,
		Payload:    notify.Payload,
		RetryCount: 0,
	})

	// Expect: получение из кеша
	mockRedis.EXPECT().
		Get(ctx, notify.ID).
		Return(notify, nil).
		Times(1)

	// Expect: ошибка при отправке
	mockSender.EXPECT().
		Send(ctx, gomock.Any()).
		Return(errors.New("send failed")).
		Times(1)

	// Expect: обновление статуса на Pending с увеличением retry count
	mockPostgres.EXPECT().
		UpdateStatus(ctx, notify.ID, domain.StatusPending, gomock.Any(), 1, gomock.Any()).
		Return(nil).
		Times(1)

	// Act
	err := consumer.Handle(ctx, payload)

	// Assert - должно вернуть nil (ошибка обработана и запланирован retry)
	if err != nil {
		t.Errorf("expected nil error for retry scenario, got %v", err)
	}
}

func TestNotifyConsumer_Handle_SendFailure_MaxRetries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks.NewMockNotifyPostgres(ctrl)
	mockRedis := mocks.NewMockNotifyRedis(ctrl)
	mockQueue := mocks.NewMockQueueProvider(ctrl)
	mockSender := mocks.NewMockSender(ctrl)

	cfg := config.NotifierConfig{
		MaxRetries: 3,
	}

	senders := map[string]domain.Sender{
		"email": mockSender,
	}

	consumer := NewNotifyConsumer(cfg, mockPostgres, mockQueue, mockRedis, senders, log.New())

	ctx := context.Background()
	notify := &domain.Notify{
		ID:         "test-id-123",
		Target:     "test@example.com",
		Channel:    "email",
		Payload:    []byte("Test message"),
		Status:     domain.StatusPending,
		RetryCount: 2, // Уже 2 попытки
	}

	payload, _ := json.Marshal(NotifyWorkerDTO{
		ID:         notify.ID,
		Target:     notify.Target,
		Channel:    notify.Channel,
		Payload:    notify.Payload,
		RetryCount: 2,
	})

	// Expect: получение из кеша
	mockRedis.EXPECT().
		Get(ctx, notify.ID).
		Return(notify, nil).
		Times(1)

	// Expect: ошибка при отправке
	mockSender.EXPECT().
		Send(ctx, gomock.Any()).
		Return(errors.New("send failed")).
		Times(1)

	// Expect: обновление статуса на Failed (достигнут лимит)
	mockPostgres.EXPECT().
		UpdateStatus(ctx, notify.ID, domain.StatusFailed, nil, 3, gomock.Any()).
		Return(nil).
		Times(1)

	// Act
	err := consumer.Handle(ctx, payload)

	// Assert
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestNotifyConsumer_Send_UnsupportedChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks.NewMockNotifyPostgres(ctrl)
	mockRedis := mocks.NewMockNotifyRedis(ctrl)
	mockQueue := mocks.NewMockQueueProvider(ctrl)

	cfg := config.NotifierConfig{MaxRetries: 3}
	senders := map[string]domain.Sender{
		"email": mocks.NewMockSender(ctrl),
	}

	consumer := NewNotifyConsumer(cfg, mockPostgres, mockQueue, mockRedis, senders, log.New())

	ctx := context.Background()
	dto := NotifyWorkerDTO{
		ID:      "test-id",
		Target:  "123456",
		Channel: "unsupported-channel", // Неподдерживаемый канал
	}

	// Act
	err := consumer.Send(ctx, dto)

	// Assert
	if err == nil {
		t.Errorf("expected error for unsupported channel, got nil")
	}
	if !errors.Is(err, errors.New("unsupported channel: unsupported-channel")) {
		// Проверяем что ошибка содержит правильное сообщение
		if err.Error() != "unsupported channel: unsupported-channel" {
			t.Errorf("expected 'unsupported channel' error, got %v", err)
		}
	}
}
