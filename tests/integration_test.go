package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adexcell/delayed-notifier/internal/notify/controller/http_router"
	"github.com/adexcell/delayed-notifier/internal/notify/domain"
	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/adexcell/delayed-notifier/internal/notify/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/wb-go/wbf/ginext"
	"github.com/adexcell/delayed-notifier/pkg/metrics"
	"github.com/adexcell/delayed-notifier/pkg/otel"
)

type MockPostgres struct {
	usecase.Postgres
	CreateNotifyFunc func(ctx context.Context, notify domain.Notify) error
}

func (m *MockPostgres) CreateNotify(ctx context.Context, notify domain.Notify) error {
	if m.CreateNotifyFunc != nil {
		return m.CreateNotifyFunc(ctx, notify)
	}
	return nil
}

func (m *MockPostgres) DeleteNotify(ctx context.Context, notifyID uuid.UUID) error {
	return nil
}

type MockRedis struct {}

func (m *MockRedis) GetNotifyStatus(ctx context.Context, key string) (domain.Status, error) { return domain.Status(0), nil }

type MockRedisWriter struct {}
func (m *MockRedisWriter) Send(ctx context.Context, task domain.NotifyStatusTask) {}

type MockRabbitWriter struct {}
func (m *MockRabbitWriter) Send(ctx context.Context, delivery dto.Delivery) {}
var routerCache *ginext.Engine
var pgMockCache *MockPostgres

func setupRouter() (*ginext.Engine, *MockPostgres) {
	if routerCache != nil {
		pgMockCache.CreateNotifyFunc = nil
		return routerCache, pgMockCache
	}
	otel.SilentModeInit()

	pgMockCache = &MockPostgres{}
	uc := usecase.New(pgMockCache, &MockRedis{}, &MockRedisWriter{}, &MockRabbitWriter{})
	
	routerCache = ginext.New("")
	m := metrics.NewHTTPServer()
	httprouter.NotifyRouter(routerCache, uc, m)
	return routerCache, pgMockCache
}

func TestIntegration_CreateNotify(t *testing.T) {
	router, _ := setupRouter()

	input := dto.CreateNotifyInput{
		RecipientEmail: "integration@example.com",
		Subject:        "Int Test",
		Body:           "Body",
		ScheduledAt:    time.Now().Add(10 * time.Hour),
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/notifies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var output dto.CreateNotifyOutput
	err := json.Unmarshal(w.Body.Bytes(), &output)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, output.ID)
}

func TestIntegration_DeleteNotify(t *testing.T) {
	router, _ := setupRouter()

	id := uuid.New().String()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/notifies/"+id, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
}
