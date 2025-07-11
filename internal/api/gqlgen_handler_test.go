package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NarthurN/CommentsSystem/internal/config"
	"github.com/NarthurN/CommentsSystem/internal/repository"
	"github.com/NarthurN/CommentsSystem/internal/service"
	"github.com/NarthurN/CommentsSystem/pkg/pubsub"
)

// Мок для Storage интерфейса
type mockStorage struct {
	repository.Storage
}

func (m *mockStorage) HealthCheck(ctx context.Context) error {
	return nil
}

func TestNewGQLGenHandler(t *testing.T) {
	mockStorage := &mockStorage{}
	ps := pubsub.New()
	svc := service.NewGQLGenService(mockStorage, ps)

	handler := NewGQLGenHandler(svc)

	if handler == nil {
		t.Fatal("NewGQLGenHandler returned nil")
	}

	if handler.service != svc {
		t.Error("Service not properly set in handler")
	}
}

func TestNewGQLGenHandlerWithConfig(t *testing.T) {
	mockStorage := &mockStorage{}
	ps := pubsub.New()
	svc := service.NewGQLGenService(mockStorage, ps)
	cfg := &config.Config{
		RequestTimeout: config.DefaultRequestTimeout,
		AllowOrigin:    config.DefaultAllowOrigin,
		AllowMethods:   config.DefaultAllowMethods,
		AllowHeaders:   config.DefaultAllowHeaders,
	}

	handler := NewGQLGenHandlerWithConfig(svc, cfg)

	if handler == nil {
		t.Fatal("NewGQLGenHandlerWithConfig returned nil")
	}

	if handler.config != cfg {
		t.Error("Config not properly set in handler")
	}
}

func TestGQLGenHandler_SetupRoutes(t *testing.T) {
	mockStorage := &mockStorage{}
	ps := pubsub.New()
	svc := service.NewGQLGenService(mockStorage, ps)
	handler := NewGQLGenHandler(svc)

	router := handler.SetupRoutes()

	if router == nil {
		t.Error("SetupRoutes returned nil router")
	}
}

func TestGQLGenHandler_HandleHealthCheck(t *testing.T) {
	mockStorage := &mockStorage{}
	ps := pubsub.New()
	svc := service.NewGQLGenService(mockStorage, ps)
	handler := NewGQLGenHandler(svc)

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	handler.HandleHealthCheck(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	expectedContentType := "application/json"
	if contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}
}

func TestGQLGenHandler_isOriginAllowed(t *testing.T) {
	tests := []struct {
		name         string
		allowOrigin  string
		testOrigin   string
		expectedBool bool
	}{
		{
			name:         "wildcard origin",
			allowOrigin:  "*",
			testOrigin:   "https://example.com",
			expectedBool: true,
		},
		{
			name:         "exact match",
			allowOrigin:  "https://example.com",
			testOrigin:   "https://example.com",
			expectedBool: true,
		},
		{
			name:         "no match",
			allowOrigin:  "https://example.com",
			testOrigin:   "https://malicious.com",
			expectedBool: false,
		},
		{
			name:         "multiple origins - match first",
			allowOrigin:  "https://example.com,https://test.com",
			testOrigin:   "https://example.com",
			expectedBool: true,
		},
		{
			name:         "multiple origins - match second",
			allowOrigin:  "https://example.com,https://test.com",
			testOrigin:   "https://test.com",
			expectedBool: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{AllowOrigin: tt.allowOrigin}
			mockStorage := &mockStorage{}
			ps := pubsub.New()
			svc := service.NewGQLGenService(mockStorage, ps)
			handler := NewGQLGenHandlerWithConfig(svc, cfg)

			result := handler.isOriginAllowed(tt.testOrigin)
			if result != tt.expectedBool {
				t.Errorf("isOriginAllowed() = %v, expected %v", result, tt.expectedBool)
			}
		})
	}
}
