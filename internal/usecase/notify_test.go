package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/internal/mocks"
	"github.com/adexcell/delayed-notifier/pkg/log"
	"go.uber.org/mock/gomock"
)

func TestNotifyUsecase_Save_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks.NewMockNotifyPostgres(ctrl)
	mockRedis := mocks.NewMockNotifyRedis(ctrl)
	mockQueue := mocks.NewMockQueueProvider(ctrl)

	usecase := New(mockPostgres, mockRedis, mockQueue, log.New())

	ctx := context.Background()
	notify := &domain.Notify{
		ID:      "test-id-123",
		Target:  "test@example.com",
		Channel: "email",
		Status:  domain.StatusPending,
	}

	// Expect: проверка на существование (не найдено)
	mockPostgres.EXPECT().
		GetNotifyByID(ctx, notify.ID).
		Return(nil, domain.ErrNotFound).
		Times(1)

	// Expect: создание записи
	mockPostgres.EXPECT().
		Create(ctx, notify).
		Return(nil).
		Times(1)

	// Act
	id, err := usecase.Save(ctx, notify)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if id != notify.ID {
		t.Errorf("expected id %s, got %s", notify.ID, id)
	}
}

func TestNotifyUsecase_Save_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks.NewMockNotifyPostgres(ctrl)
	mockRedis := mocks.NewMockNotifyRedis(ctrl)
	mockQueue := mocks.NewMockQueueProvider(ctrl)

	usecase := New(mockPostgres, mockRedis, mockQueue, log.New())

	ctx := context.Background()
	notify := &domain.Notify{
		ID:      "test-id-123",
		Target:  "test@example.com",
		Channel: "email",
	}

	existingNotify := &domain.Notify{ID: notify.ID}

	// Expect: notify уже существует
	mockPostgres.EXPECT().
		GetNotifyByID(ctx, notify.ID).
		Return(existingNotify, nil).
		Times(1)

	// Act
	id, err := usecase.Save(ctx, notify)

	// Assert
	if !errors.Is(err, domain.ErrNotifyAlreadyExists) {
		t.Errorf("expected ErrNotifyAlreadyExists, got %v", err)
	}
	if id != notify.ID {
		t.Errorf("expected id %s, got %s", notify.ID, id)
	}
}

func TestNotifyUsecase_GetByID_FromCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks.NewMockNotifyPostgres(ctrl)
	mockRedis := mocks.NewMockNotifyRedis(ctrl)
	mockQueue := mocks.NewMockQueueProvider(ctrl)

	usecase := New(mockPostgres, mockRedis, mockQueue, log.New())

	ctx := context.Background()
	expectedNotify := &domain.Notify{
		ID:      "test-id-123",
		Target:  "test@example.com",
		Channel: "email",
	}

	// Expect: получение из Redis кеша
	mockRedis.EXPECT().
		Get(ctx, expectedNotify.ID).
		Return(expectedNotify, nil).
		Times(1)

	// Act
	notify, err := usecase.GetByID(ctx, expectedNotify.ID)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if notify.ID != expectedNotify.ID {
		t.Errorf("expected notify ID %s, got %s", expectedNotify.ID, notify.ID)
	}
}

func TestNotifyUsecase_GetByID_FromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks.NewMockNotifyPostgres(ctrl)
	mockRedis := mocks.NewMockNotifyRedis(ctrl)
	mockQueue := mocks.NewMockQueueProvider(ctrl)

	usecase := New(mockPostgres, mockRedis, mockQueue, log.New())

	ctx := context.Background()
	expectedNotify := &domain.Notify{
		ID:      "test-id-123",
		Target:  "test@example.com",
		Channel: "email",
	}

	// Expect: кеш промах
	mockRedis.EXPECT().
		Get(ctx, expectedNotify.ID).
		Return(nil, errors.New("redis: nil")).
		Times(1)

	// Expect: получение из БД
	mockPostgres.EXPECT().
		GetNotifyByID(ctx, expectedNotify.ID).
		Return(expectedNotify, nil).
		Times(1)

	// Expect: сохранение в кеш
	mockRedis.EXPECT().
		SetWithExpiration(ctx, expectedNotify).
		Return(nil).
		Times(1)

	// Act
	notify, err := usecase.GetByID(ctx, expectedNotify.ID)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if notify.ID != expectedNotify.ID {
		t.Errorf("expected notify ID %s, got %s", expectedNotify.ID, notify.ID)
	}
}

func TestNotifyUsecase_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks.NewMockNotifyPostgres(ctrl)
	mockRedis := mocks.NewMockNotifyRedis(ctrl)
	mockQueue := mocks.NewMockQueueProvider(ctrl)

	usecase := New(mockPostgres, mockRedis, mockQueue, log.New())

	ctx := context.Background()
	notifyID := "non-existent-id"

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
	notify, err := usecase.GetByID(ctx, notifyID)

	// Assert
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	if notify != nil {
		t.Errorf("expected nil notify, got %v", notify)
	}
}

func TestNotifyUsecase_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks.NewMockNotifyPostgres(ctrl)
	mockRedis := mocks.NewMockNotifyRedis(ctrl)
	mockQueue := mocks.NewMockQueueProvider(ctrl)

	usecase := New(mockPostgres, mockRedis, mockQueue, log.New())

	ctx := context.Background()
	notifyID := "test-id-123"

	// Expect: удаление из БД
	mockPostgres.EXPECT().
		DeleteByID(ctx, notifyID).
		Return(nil).
		Times(1)

	// Act
	err := usecase.Delete(ctx, notifyID)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestNotifyUsecase_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks.NewMockNotifyPostgres(ctrl)
	mockRedis := mocks.NewMockNotifyRedis(ctrl)
	mockQueue := mocks.NewMockQueueProvider(ctrl)

	usecase := New(mockPostgres, mockRedis, mockQueue, log.New())

	ctx := context.Background()
	limit := 10
	offset := 0

	expectedNotifies := []*domain.Notify{
		{ID: "id-1", Target: "user1@example.com"},
		{ID: "id-2", Target: "user2@example.com"},
	}

	// Expect: получение списка из БД
	mockPostgres.EXPECT().
		List(ctx, limit, offset).
		Return(expectedNotifies, nil).
		Times(1)

	// Act
	notifies, err := usecase.List(ctx, limit, offset)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(notifies) != len(expectedNotifies) {
		t.Errorf("expected %d notifies, got %d", len(expectedNotifies), len(notifies))
	}
}
