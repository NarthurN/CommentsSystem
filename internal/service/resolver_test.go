package service

import (
	"testing"

	"github.com/NarthurN/CommentsSystem/internal/repository"
	"github.com/NarthurN/CommentsSystem/pkg/pubsub"
)

// Мок для Storage интерфейса
type mockStorage struct {
	repository.Storage
}

func TestNewResolver(t *testing.T) {
	mockStorage := &mockStorage{}
	ps := pubsub.New()

	resolver := NewResolver(mockStorage, ps)

	if resolver == nil {
		t.Fatal("NewResolver returned nil")
	}

	if resolver.storage != mockStorage {
		t.Error("Storage not properly set in resolver")
	}

	if resolver.pubsub != ps {
		t.Error("PubSub not properly set in resolver")
	}
}
