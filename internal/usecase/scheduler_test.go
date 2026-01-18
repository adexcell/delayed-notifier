package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/adexcell/delayed-notifier/config"
	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/internal/mocks"
	"github.com/adexcell/delayed-notifier/pkg/log"
	"go.uber.org/mock/gomock"
)

func TestScheduler_Process_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks.NewMockNotifyPostgres(ctrl)
	mockQueue := mocks.NewMockQueueProvider(ctrl)

	cfg := config.NotifierConfig{
		BatchSize:         10,
		VisibilityTimeout: 60 * time.Second,
		MaxRetries:        3,
	}

	scheduler := NewScheduler(mockPostgres, mockQueue, cfg, log.New())

	// Используем приватный метод process для теста, чтобы не запускать бесконечный цикл Run
	// Но так как process приватный, мы не можем его вызвать из update_test.go если он в другом пакете.
	// Поэтому тест должен быть в package usecase.

	ctx := context.Background()

	notifies := []*domain.Notify{
		{ID: "1", Status: domain.StatusPending, ScheduledAt: time.Now()},
		{ID: "2", Status: domain.StatusPending, ScheduledAt: time.Now()},
	}

	// Expect: LockAndFetchReady возвращает уведомления
	mockPostgres.EXPECT().
		LockAndFetchReady(ctx, cfg.BatchSize, cfg.VisibilityTimeout).
		Return(notifies, nil).
		Times(1)

	// Expect: Publish для каждого уведомления
	// Мы не можем гарантировать порядок, если горутины, но тут синхронно.
	mockQueue.EXPECT().Publish(ctx, notifies[0]).Return(nil).Times(1)
	mockQueue.EXPECT().Publish(ctx, notifies[1]).Return(nil).Times(1)

	// Для теста экспортируем метод или используем reflect/linkname, но проще просто проверить логику
	// вызова зависимостей.
	// В Go принято тестировать публичный API. Если process приватный, тестируем через Run?
	// Run запускает тикер. Это долго.
	// Правильнее было бы сделать process публичным (Process) или тестировать внутреннюю логику.
	// В данном случае, так как мы пишем тесты в package usecase (а не usecase_test), мы имеем доступ к приватным методам.
	// (scheduler.go находится в package usecase)

	// Приведение типа к конкретной структуре чтобы вызвать приватный метод, если интерфейс скрывает его?
	// NewScheduler возвращает domain.Scheduler интерфейс, у которого только Run.
	// Нам нужно скастить к структуре *Scheduler.

	s, ok := scheduler.(*Scheduler)
	if !ok {
		t.Fatal("expected *Scheduler type")
	}

	s.process(ctx)
}

func TestScheduler_Process_PublishError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks.NewMockNotifyPostgres(ctrl)
	mockQueue := mocks.NewMockQueueProvider(ctrl)

	cfg := config.NotifierConfig{MaxRetries: 3}
	scheduler := NewScheduler(mockPostgres, mockQueue, cfg, log.New())
	s := scheduler.(*Scheduler)

	ctx := context.Background()
	notify := &domain.Notify{ID: "1", RetryCount: 0}

	mockPostgres.EXPECT().
		LockAndFetchReady(ctx, cfg.BatchSize, cfg.VisibilityTimeout).
		Return([]*domain.Notify{notify}, nil).
		Times(1)

	// Expect: ошибка публикации
	mockQueue.EXPECT().
		Publish(ctx, notify).
		Return(errors.New("queue error")).
		Times(1)

	// Expect: обновление статуса (retry count + 1)
	mockPostgres.EXPECT().
		UpdateStatus(ctx, notify.ID, domain.StatusPending, gomock.Any(), 1, gomock.Any()).
		Return(nil).
		Times(1)

	s.process(ctx)
}

func TestScheduler_Process_MaxRetries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks.NewMockNotifyPostgres(ctrl)
	mockQueue := mocks.NewMockQueueProvider(ctrl)

	cfg := config.NotifierConfig{MaxRetries: 3}
	scheduler := NewScheduler(mockPostgres, mockQueue, cfg, log.New())
	s := scheduler.(*Scheduler)

	ctx := context.Background()
	notify := &domain.Notify{ID: "1", RetryCount: 3} // Уже 3 попытки (достигнут максимум при следующей ошибке?)
	// Логика в scheduler.go:
	// if n.RetryCount < s.maxRetries { status = Pending } else { status = Failed }
	// Если MaxRetries=3, и RetryCount=3. 3 < 3 is false. -> Failed.

	mockPostgres.EXPECT().
		LockAndFetchReady(ctx, cfg.BatchSize, cfg.VisibilityTimeout).
		Return([]*domain.Notify{notify}, nil).
		Times(1)

	mockQueue.EXPECT().
		Publish(ctx, notify).
		Return(errors.New("queue error")).
		Times(1)

	// Expect: статус Failed
	mockPostgres.EXPECT().
		UpdateStatus(ctx, notify.ID, domain.StatusFailed, gomock.Any(), 0, gomock.Any()).
		Return(nil).
		Times(1)

	s.process(ctx)
}
