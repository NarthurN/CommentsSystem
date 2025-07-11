package service

import (
	"context"
	"testing"

	"github.com/NarthurN/CommentsSystem/internal/config"
	"github.com/NarthurN/CommentsSystem/internal/repository"
	"github.com/NarthurN/CommentsSystem/pkg/pubsub"
)

func TestNewGQLGenService(t *testing.T) {
	mockStorage := &mockStorage{}
	ps := pubsub.New()

	service := NewGQLGenService(mockStorage, ps)

	if service == nil {
		t.Fatal("NewGQLGenService returned nil")
	}

	if service.storage != mockStorage {
		t.Error("Storage not properly set in service")
	}

	if service.pubsub != ps {
		t.Error("PubSub not properly set in service")
	}

	if service.server == nil {
		t.Error("Server should be initialized")
	}

	if service.resolver == nil {
		t.Error("Resolver should be initialized")
	}
}

func TestNewGQLGenServiceWithConfig(t *testing.T) {
	mockStorage := &mockStorage{}
	ps := pubsub.New()
	cfg := &config.Config{
		KeepAlivePing:       config.DefaultKeepAlivePing,
		AllowOrigin:         config.DefaultAllowOrigin,
		PlaygroundTitle:     config.DefaultPlaygroundTitle,
		GraphQLEndpoint:     config.DefaultGraphQLEndpoint,
		EnableIntrospection: true,
	}

	service := NewGQLGenServiceWithConfig(mockStorage, ps, cfg)

	if service == nil {
		t.Fatal("NewGQLGenServiceWithConfig returned nil")
	}

	if service.config != cfg {
		t.Error("Config not properly set in service")
	}
}

func TestGQLGenService_GetHandler(t *testing.T) {
	mockStorage := &mockStorage{}
	ps := pubsub.New()
	service := NewGQLGenService(mockStorage, ps)

	handler := service.GetHandler()

	if handler == nil {
		t.Error("GetHandler returned nil")
	}
}

func TestGQLGenService_GetPlaygroundHandler(t *testing.T) {
	mockStorage := &mockStorage{}
	ps := pubsub.New()
	service := NewGQLGenService(mockStorage, ps)

	handler := service.GetPlaygroundHandler()

	if handler == nil {
		t.Error("GetPlaygroundHandler returned nil")
	}
}

func TestGQLGenService_HealthCheck(t *testing.T) {
	mockStorage := &mockHealthCheckStorage{}
	ps := pubsub.New()
	service := NewGQLGenService(mockStorage, ps)

	err := service.HealthCheck(context.Background())

	if err != nil {
		t.Errorf("HealthCheck returned error: %v", err)
	}
}

func TestGQLGenService_GetConfig(t *testing.T) {
	mockStorage := &mockStorage{}
	ps := pubsub.New()
	cfg := &config.Config{
		KeepAlivePing:   config.DefaultKeepAlivePing,
		AllowOrigin:     config.DefaultAllowOrigin,
		PlaygroundTitle: config.DefaultPlaygroundTitle,
		GraphQLEndpoint: config.DefaultGraphQLEndpoint,
	}

	service := NewGQLGenServiceWithConfig(mockStorage, ps, cfg)

	returnedCfg := service.GetConfig()

	if returnedCfg != cfg {
		t.Error("GetConfig returned wrong config")
	}
}

func TestGQLGenService_GetSubscribersCount(t *testing.T) {
	mockStorage := &mockStorage{}
	ps := pubsub.New()
	service := NewGQLGenService(mockStorage, ps)

	count := service.GetSubscribersCount()

	// Поскольку GetTotalSubscribers удален, ожидаем 0
	if count != 0 {
		t.Errorf("Expected 0 subscribers, got %d", count)
	}
}

// Мок с поддержкой HealthCheck
type mockHealthCheckStorage struct {
	repository.Storage
}

func (m *mockHealthCheckStorage) HealthCheck(ctx context.Context) error {
	return nil
}
