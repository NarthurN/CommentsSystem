package repository_test

import (
	"context"
	"testing"

	"github.com/NarthurN/CommentsSystem/internal/model"
	"github.com/NarthurN/CommentsSystem/internal/repository"
	"github.com/google/uuid"
)

// TestPostgresStorage_CreatePost тестирует создание поста в PostgreSQL
func TestPostgresStorage_CreatePost(t *testing.T) {
	// Проверяем, что тест запущен с флагом интеграционных тестов
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Для интеграционного теста нужна реальная база данных
	// В CI/CD можно использовать testcontainers или Docker
	dsn := "postgres://user:password@localhost:5432/postsdb_test?sslmode=disable"

	ctx := context.Background()
	storage, err := repository.NewPostgresStorage(ctx, dsn)
	if err != nil {
		t.Skipf("Failed to connect to test database: %v", err)
	}
	defer storage.Close()

	// Создаем тестовый пост
	post := &model.Post{
		Title:   "Test Post",
		Content: "This is a test post content",
	}

	// Создаем пост
	createdPost, err := storage.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}

	// Проверяем, что пост создан
	if createdPost.ID == uuid.Nil {
		t.Error("Expected post to have ID")
	}
	if createdPost.Title != post.Title {
		t.Errorf("Expected title %q, got %q", post.Title, createdPost.Title)
	}
	if createdPost.Content != post.Content {
		t.Errorf("Expected content %q, got %q", post.Content, createdPost.Content)
	}
	if !createdPost.CommentsEnabled {
		t.Error("Expected comments to be enabled by default")
	}
	if createdPost.CreatedAt.IsZero() {
		t.Error("Expected created_at to be set")
	}

	// Проверяем, что можем получить пост по ID
	fetchedPost, err := storage.GetPost(ctx, createdPost.ID)
	if err != nil {
		t.Fatalf("Failed to get post by ID: %v", err)
	}
	if fetchedPost == nil {
		t.Fatal("Expected to find post")
	}
	if fetchedPost.ID != createdPost.ID {
		t.Errorf("Expected ID %v, got %v", createdPost.ID, fetchedPost.ID)
	}
}

// TestPostgresStorage_CommentHierarchy тестирует иерархическую структуру комментариев
func TestPostgresStorage_CommentHierarchy(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	dsn := "postgres://user:password@localhost:5432/postsdb_test?sslmode=disable"
	ctx := context.Background()
	storage, err := repository.NewPostgresStorage(ctx, dsn)
	if err != nil {
		t.Skipf("Failed to connect to test database: %v", err)
	}
	defer storage.Close()

	// Создаем пост
	post := &model.Post{
		Title:   "Post with Comments",
		Content: "This post will have comments",
	}
	createdPost, err := storage.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}

	// Создаем корневой комментарий
	rootComment := &model.Comment{
		PostID:  createdPost.ID,
		Content: "This is a root comment",
	}
	createdRootComment, err := storage.CreateComment(ctx, rootComment)
	if err != nil {
		t.Fatalf("Failed to create root comment: %v", err)
	}

	// Создаем дочерний комментарий
	childComment := &model.Comment{
		PostID:   createdPost.ID,
		ParentID: &createdRootComment.ID,
		Content:  "This is a child comment",
	}
	createdChildComment, err := storage.CreateComment(ctx, childComment)
	if err != nil {
		t.Fatalf("Failed to create child comment: %v", err)
	}

	// Получаем все комментарии для поста
	comments, err := storage.GetCommentsByPostID(ctx, createdPost.ID)
	if err != nil {
		t.Fatalf("Failed to get comments: %v", err)
	}

	// Проверяем, что получили оба комментария
	if len(comments) != 2 {
		t.Errorf("Expected 2 comments, got %d", len(comments))
	}

	// Проверяем иерархию
	foundRoot := false
	foundChild := false
	for _, comment := range comments {
		if comment.ID == createdRootComment.ID {
			foundRoot = true
			if comment.ParentID != nil {
				t.Error("Root comment should not have parent")
			}
		}
		if comment.ID == createdChildComment.ID {
			foundChild = true
			if comment.ParentID == nil || *comment.ParentID != createdRootComment.ID {
				t.Error("Child comment should have correct parent")
			}
		}
	}

	if !foundRoot {
		t.Error("Root comment not found")
	}
	if !foundChild {
		t.Error("Child comment not found")
	}
}
