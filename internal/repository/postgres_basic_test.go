package repository

import (
	"testing"

	"github.com/NarthurN/CommentsSystem/internal/repository/converter"
)

func TestNewPostgresStorage_ConvertersInitialization(t *testing.T) {
	// Тестируем только инициализацию конвертеров без реального подключения к БД
	postConverter := converter.NewPostConverter()
	commentConverter := converter.NewCommentConverter()
	treeConverter := converter.NewTreeConverter()

	if postConverter == nil {
		t.Error("PostConverter should be initialized")
	}

	if commentConverter == nil {
		t.Error("CommentConverter should be initialized")
	}

	if treeConverter == nil {
		t.Error("TreeConverter should be initialized")
	}
}
