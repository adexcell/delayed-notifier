package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/internal/mocks"
	"github.com/adexcell/delayed-notifier/pkg/log"
	"github.com/adexcell/delayed-notifier/pkg/router"
	"go.uber.org/mock/gomock"
)

func TestNotifyHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockNotifyUsecase(ctrl)

	r := router.New(router.Config{GinMode: "test"})
	handler := NewNotifyHandler(mockUsecase, log.New())
	handler.Register(r)

	futureTime := time.Now().Add(10 * time.Minute)
	requestBody := NotifyControllerDTO{
		Payload:     json.RawMessage(`"test message"`),
		Target:      "test@example.com",
		Channel:     "email",
		ScheduledAt: futureTime,
	}

	body, _ := json.Marshal(requestBody)

	// Expect: успешное сохранение
	mockUsecase.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return("generated-id", nil).
		Times(1)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/notify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["id"] == "" {
		t.Errorf("expected id in response, got empty")
	}
}

func TestNotifyHandler_Create_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockNotifyUsecase(ctrl)

	r := router.New(router.Config{GinMode: "test"})
	handler := NewNotifyHandler(mockUsecase, log.New())
	handler.Register(r)

	invalidJSON := []byte(`{invalid json}`)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/notify", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestNotifyHandler_Create_PastScheduledAt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockNotifyUsecase(ctrl)

	r := router.New(router.Config{GinMode: "test"})
	handler := NewNotifyHandler(mockUsecase, log.New())
	handler.Register(r)

	pastTime := time.Now().Add(-10 * time.Minute)
	requestBody := NotifyControllerDTO{
		Payload:     json.RawMessage(`"test message"`),
		Target:      "test@example.com",
		Channel:     "email",
		ScheduledAt: pastTime,
	}

	body, _ := json.Marshal(requestBody)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/notify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestNotifyHandler_Create_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockNotifyUsecase(ctrl)

	r := router.New(router.Config{GinMode: "test"})
	handler := NewNotifyHandler(mockUsecase, log.New())
	handler.Register(r)

	futureTime := time.Now().Add(10 * time.Minute)
	requestBody := NotifyControllerDTO{
		Payload:     json.RawMessage(`"test message"`),
		Target:      "test@example.com",
		Channel:     "email",
		ScheduledAt: futureTime,
	}

	body, _ := json.Marshal(requestBody)

	// Expect: notify уже существует
	mockUsecase.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return("", domain.ErrNotifyAlreadyExists).
		Times(1)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/notify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, w.Code)
	}
}

func TestNotifyHandler_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockNotifyUsecase(ctrl)

	r := router.New(router.Config{GinMode: "test"})
	handler := NewNotifyHandler(mockUsecase, log.New())
	handler.Register(r)

	notifyID := "550e8400-e29b-41d4-a716-446655440000"
	expectedNotify := &domain.Notify{
		ID:          notifyID,
		Target:      "test@example.com",
		Channel:     "email",
		Status:      domain.StatusPending,
		ScheduledAt: time.Now(),
		CreatedAt:   time.Now(),
		RetryCount:  0,
	}

	// Expect: получение notify
	mockUsecase.EXPECT().
		GetByID(gomock.Any(), notifyID).
		Return(expectedNotify, nil).
		Times(1)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/notify/"+notifyID, nil)
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response NotifyResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.ID != notifyID {
		t.Errorf("expected notify ID %s, got %s", notifyID, response.ID)
	}
}

func TestNotifyHandler_Get_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockNotifyUsecase(ctrl)

	r := router.New(router.Config{GinMode: "test"})
	handler := NewNotifyHandler(mockUsecase, log.New())
	handler.Register(r)

	invalidID := "not-a-uuid"

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/notify/"+invalidID, nil)
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestNotifyHandler_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockNotifyUsecase(ctrl)

	r := router.New(router.Config{GinMode: "test"})
	handler := NewNotifyHandler(mockUsecase, log.New())
	handler.Register(r)

	notifyID := "550e8400-e29b-41d4-a716-446655440000"

	// Expect: notify не найдено
	mockUsecase.EXPECT().
		GetByID(gomock.Any(), notifyID).
		Return(nil, domain.ErrNotFound).
		Times(1)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/notify/"+notifyID, nil)
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestNotifyHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockNotifyUsecase(ctrl)

	r := router.New(router.Config{GinMode: "test"})
	handler := NewNotifyHandler(mockUsecase, log.New())
	handler.Register(r)

	notifyID := "550e8400-e29b-41d4-a716-446655440000"

	// Expect: успешное удаление
	mockUsecase.EXPECT().
		Delete(gomock.Any(), notifyID).
		Return(nil).
		Times(1)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/notify/"+notifyID, nil)
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestNotifyHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockNotifyUsecase(ctrl)

	r := router.New(router.Config{GinMode: "test"})
	handler := NewNotifyHandler(mockUsecase, log.New())
	handler.Register(r)

	notifyID := "550e8400-e29b-41d4-a716-446655440000"

	// Expect: notify не найдено (но это ОК для DELETE - идемпотентность)
	mockUsecase.EXPECT().
		Delete(gomock.Any(), notifyID).
		Return(domain.ErrNotFound).
		Times(1)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/notify/"+notifyID, nil)
	r.ServeHTTP(w, req)

	// Assert - логика в коде пропускает ErrNotFound
	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestNotifyHandler_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockNotifyUsecase(ctrl)

	r := router.New(router.Config{GinMode: "test"})
	handler := NewNotifyHandler(mockUsecase, log.New())
	handler.Register(r)

	expectedNotifies := []*domain.Notify{
		{ID: "id-1", Target: "user1@example.com", Status: domain.StatusPending},
		{ID: "id-2", Target: "user2@example.com", Status: domain.StatusSent},
	}

	// Expect: получение списка
	mockUsecase.EXPECT().
		List(gomock.Any(), 50, 0). // Default limit=50, offset=0
		Return(expectedNotifies, nil).
		Times(1)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/notify", nil)
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []*domain.Notify
	json.Unmarshal(w.Body.Bytes(), &response)
	if len(response) != len(expectedNotifies) {
		t.Errorf("expected %d notifies, got %d", len(expectedNotifies), len(response))
	}
}

func TestNotifyHandler_List_WithPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockNotifyUsecase(ctrl)

	r := router.New(router.Config{GinMode: "test"})
	handler := NewNotifyHandler(mockUsecase, log.New())
	handler.Register(r)

	expectedNotifies := []*domain.Notify{
		{ID: "id-3", Target: "user3@example.com"},
	}

	// Expect: получение списка с параметрами
	mockUsecase.EXPECT().
		List(gomock.Any(), 10, 20).
		Return(expectedNotifies, nil).
		Times(1)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/notify?limit=10&offset=20", nil)
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestNotifyHandler_List_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockNotifyUsecase(ctrl)

	r := router.New(router.Config{GinMode: "test"})
	handler := NewNotifyHandler(mockUsecase, log.New())
	handler.Register(r)

	// Expect: ошибка БД
	mockUsecase.EXPECT().
		List(gomock.Any(), 50, 0).
		Return(nil, errors.New("database error")).
		Times(1)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/notify", nil)
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
