package service

import (
	"github.com/NarthurN/CommentsSystem/internal/repository"
	"github.com/NarthurN/CommentsSystem/pkg/pubsub"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	storage repository.Storage
	pubsub  *pubsub.PubSub
}

// NewResolver создает новый экземпляр Resolver с зависимостями
func NewResolver(storage repository.Storage, ps *pubsub.PubSub) *Resolver {
	return &Resolver{
		storage: storage,
		pubsub:  ps,
	}
}
