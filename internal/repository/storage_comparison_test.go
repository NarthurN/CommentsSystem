package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/NarthurN/CommentsSystem/internal/model"
	"github.com/NarthurN/CommentsSystem/internal/repository"
	"github.com/google/uuid"
)

// TestStorageComparison демонстрирует идентичное поведение PostgreSQL и In-Memory хранилищ
func TestStorageComparison(t *testing.T) {
	ctx := context.Background()

	// Создаем оба типа хранилища
	memoryStorage := repository.NewMemoryStorage()
	defer memoryStorage.Close()

	// PostgreSQL storage будет пропущен если не доступен
	var postgresStorage repository.Storage
	postgresStorage, err := repository.NewPostgresStorage(ctx, "postgres://user:password@localhost:5432/postsdb_test?sslmode=disable")
	if err != nil {
		t.Logf("PostgreSQL not available, testing only in-memory: %v", err)
		postgresStorage = nil
	} else {
		defer postgresStorage.Close()
	}

	storages := []struct {
		name    string
		storage repository.Storage
	}{
		{"Memory", memoryStorage},
	}

	if postgresStorage != nil {
		storages = append(storages, struct {
			name    string
			storage repository.Storage
		}{"PostgreSQL", postgresStorage})
	}

	// Тестируем идентичное поведение для всех хранилищ
	for _, s := range storages {
		t.Run(s.name, func(t *testing.T) {
			testStorageBehavior(t, s.storage)
		})
	}
}

// testStorageBehavior выполняет идентичные тесты для любого типа хранилища
func testStorageBehavior(t *testing.T, storage repository.Storage) {
	ctx := context.Background()

	// Тест 1: Создание поста
	post := &model.Post{
		Title:   "Comparison Test Post",
		Content: "Testing storage comparison",
	}
	createdPost, err := storage.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}

	// Проверяем базовые свойства
	if createdPost.ID == uuid.Nil {
		t.Error("Expected post to have ID")
	}
	if !createdPost.CommentsEnabled {
		t.Error("Expected comments to be enabled by default")
	}

	// Тест 2: Создание иерархии комментариев
	rootComment, err := storage.CreateComment(ctx, &model.Comment{
		PostID:  createdPost.ID,
		Content: "Root comment",
	})
	if err != nil {
		t.Fatalf("Failed to create root comment: %v", err)
	}

	childComment, err := storage.CreateComment(ctx, &model.Comment{
		PostID:   createdPost.ID,
		ParentID: &rootComment.ID,
		Content:  "Child comment",
	})
	if err != nil {
		t.Fatalf("Failed to create child comment: %v", err)
	}

	_, err = storage.CreateComment(ctx, &model.Comment{
		PostID:   createdPost.ID,
		ParentID: &childComment.ID,
		Content:  "Grandchild comment",
	})
	if err != nil {
		t.Fatalf("Failed to create grandchild comment: %v", err)
	}

	// Тест 3: Получение плоского списка комментариев
	comments, err := storage.GetCommentsByPostID(ctx, createdPost.ID)
	if err != nil {
		t.Fatalf("Failed to get comments: %v", err)
	}
	if len(comments) != 3 {
		t.Errorf("Expected 3 comments, got %d", len(comments))
	}

	// Тест 4: Получение иерархии комментариев
	tree, err := storage.GetCommentTree(ctx, createdPost.ID)
	if err != nil {
		t.Fatalf("Failed to get comment tree: %v", err)
	}

	// Проверяем структуру дерева
	if len(tree) != 1 {
		t.Errorf("Expected 1 root comment, got %d", len(tree))
	}

	// Проверяем глубину вложенности: root -> child -> grandchild
	rootNode := tree[0]
	if rootNode.Comment.Content != "Root comment" {
		t.Error("Root comment mismatch")
	}
	if len(rootNode.Children) != 1 {
		t.Error("Expected 1 child")
	}

	childNode := rootNode.Children[0]
	if childNode.Comment.Content != "Child comment" {
		t.Error("Child comment mismatch")
	}
	if len(childNode.Children) != 1 {
		t.Error("Expected 1 grandchild")
	}

	grandchildNode := childNode.Children[0]
	if grandchildNode.Comment.Content != "Grandchild comment" {
		t.Error("Grandchild comment mismatch")
	}
	if len(grandchildNode.Children) != 0 {
		t.Error("Expected no great-grandchildren")
	}

	// Тест 5: Отключение комментариев
	err = storage.TogglePostComments(ctx, createdPost.ID, false)
	if err != nil {
		t.Fatalf("Failed to disable comments: %v", err)
	}

	// Проверяем, что нельзя создать комментарий к посту с отключенными комментариями
	_, err = storage.CreateComment(ctx, &model.Comment{
		PostID:  createdPost.ID,
		Content: "Should fail",
	})
	if err == nil {
		t.Error("Expected error when creating comment for post with disabled comments")
	}

	// Тест 6: Каскадное удаление
	err = storage.DeleteComment(ctx, rootComment.ID)
	if err != nil {
		t.Fatalf("Failed to delete root comment: %v", err)
	}

	// Проверяем, что все комментарии удалены
	err = storage.TogglePostComments(ctx, createdPost.ID, true) // Включаем обратно для проверки
	if err != nil {
		t.Fatalf("Failed to enable comments: %v", err)
	}

	remainingComments, err := storage.GetCommentsByPostID(ctx, createdPost.ID)
	if err != nil {
		t.Fatalf("Failed to get remaining comments: %v", err)
	}
	if len(remainingComments) != 0 {
		t.Errorf("Expected 0 comments after cascade delete, got %d", len(remainingComments))
	}

	// Тест 7: Health check
	err = storage.HealthCheck(ctx)
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}

	t.Logf("✅ Storage '%T' passed all compatibility tests", storage)
}

// BenchmarkStorageComparison сравнивает производительность разных типов хранилища
func BenchmarkStorageComparison(b *testing.B) {
	ctx := context.Background()

	// Benchmark для In-Memory
	b.Run("Memory", func(b *testing.B) {
		storage := repository.NewMemoryStorage()
		defer storage.Close()
		benchmarkStorageOperations(b, storage)
	})

	// Benchmark для PostgreSQL (если доступен)
	postgresStorage, err := repository.NewPostgresStorage(ctx, "postgres://user:password@localhost:5432/postsdb_test?sslmode=disable")
	if err == nil {
		defer postgresStorage.Close()
		b.Run("PostgreSQL", func(b *testing.B) {
			benchmarkStorageOperations(b, postgresStorage)
		})
	}
}

// benchmarkStorageOperations выполняет бенчмарк операций для любого хранилища
func benchmarkStorageOperations(b *testing.B, storage repository.Storage) {
	ctx := context.Background()

	// Создаем один пост для всех операций
	post := &model.Post{
		Title:   "Benchmark Post",
		Content: "Benchmark content",
	}
	createdPost, err := storage.CreatePost(ctx, post)
	if err != nil {
		b.Fatalf("Failed to create post: %v", err)
	}

	b.ResetTimer()

	// Бенчмарк создания комментариев
	for i := 0; i < b.N; i++ {
		comment := &model.Comment{
			PostID:  createdPost.ID,
			Content: "Benchmark comment",
		}
		_, err := storage.CreateComment(ctx, comment)
		if err != nil {
			b.Fatalf("Failed to create comment: %v", err)
		}
	}
}

// TestStorageInterface проверяет, что оба хранилища реализуют интерфейс Storage
func TestStorageInterface(t *testing.T) {
	var _ repository.Storage = repository.NewMemoryStorage()

	ctx := context.Background()
	postgresStorage, err := repository.NewPostgresStorage(ctx, "dummy_dsn")
	if err == nil { // Если создание прошло успешно (обычно упадет на подключении)
		var _ repository.Storage = postgresStorage
		postgresStorage.Close()
	}

	t.Log("✅ Both storage types implement repository.Storage interface")
}

// TestStorageFeatures демонстрирует ключевые особенности каждого типа хранилища
func TestStorageFeatures(t *testing.T) {
	t.Run("MemoryStorageFeatures", func(t *testing.T) {
		storage := repository.NewMemoryStorage()
		defer storage.Close()

		// Особенность 1: Мгновенный запуск
		start := time.Now()
		ctx := context.Background()
		err := storage.HealthCheck(ctx)
		duration := time.Since(start)

		if err != nil {
			t.Errorf("Health check failed: %v", err)
		}
		if duration > time.Millisecond {
			t.Errorf("Expected instant startup, took %v", duration)
		}

		// Особенность 2: Thread-safety
		const numGoroutines = 10
		const opsPerGoroutine = 5

		// Создаем пост для тестирования
		post, err := storage.CreatePost(ctx, &model.Post{
			Title:   "Thread Safety Test",
			Content: "Testing concurrent access",
		})
		if err != nil {
			t.Fatalf("Failed to create post: %v", err)
		}

		// Запускаем concurrent операции
		done := make(chan bool, numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				for j := 0; j < opsPerGoroutine; j++ {
					storage.CreateComment(ctx, &model.Comment{
						PostID:  post.ID,
						Content: "Concurrent comment",
					})
				}
				done <- true
			}(i)
		}

		// Ждем завершения
		for i := 0; i < numGoroutines; i++ {
			<-done
		}

		// Проверяем результат
		comments, err := storage.GetCommentsByPostID(ctx, post.ID)
		if err != nil {
			t.Errorf("Failed to get comments: %v", err)
		}

		expected := numGoroutines * opsPerGoroutine
		if len(comments) != expected {
			t.Errorf("Expected %d comments, got %d (thread-safety issue)", expected, len(comments))
		}

		t.Log("✅ In-Memory storage demonstrates instant startup and thread-safety")
	})
}
