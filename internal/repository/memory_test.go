package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/NarthurN/CommentsSystem/internal/model"
	"github.com/NarthurN/CommentsSystem/internal/repository"
	"github.com/google/uuid"
)

// TestMemoryStorage_CreatePost тестирует создание поста в памяти
func TestMemoryStorage_CreatePost(t *testing.T) {
	storage := repository.NewMemoryStorage()
	defer storage.Close()

	ctx := context.Background()

	// Создаем тестовый пост
	post := &model.Post{
		Title:   "Test Memory Post",
		Content: "This is a test post in memory",
	}

	// Создаем пост
	createdPost, err := storage.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}

	// Проверяем, что пост создан правильно
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

// TestMemoryStorage_GetPosts тестирует получение списка постов с пагинацией
func TestMemoryStorage_GetPosts(t *testing.T) {
	storage := repository.NewMemoryStorage()
	defer storage.Close()

	ctx := context.Background()

	// Создаем несколько постов
	posts := []*model.Post{
		{Title: "Post 1", Content: "Content 1"},
		{Title: "Post 2", Content: "Content 2"},
		{Title: "Post 3", Content: "Content 3"},
	}

	var createdPosts []*model.Post
	for _, post := range posts {
		created, err := storage.CreatePost(ctx, post)
		if err != nil {
			t.Fatalf("Failed to create post: %v", err)
		}
		createdPosts = append(createdPosts, created)
		time.Sleep(time.Millisecond) // Обеспечиваем разное время создания
	}

	// Тестируем получение всех постов
	allPosts, err := storage.GetPosts(ctx, 10, 0)
	if err != nil {
		t.Fatalf("Failed to get posts: %v", err)
	}
	if len(allPosts) != 3 {
		t.Errorf("Expected 3 posts, got %d", len(allPosts))
	}

	// Проверяем сортировку (новые первыми)
	if !allPosts[0].CreatedAt.After(allPosts[1].CreatedAt) {
		t.Error("Expected posts to be sorted by creation time (newest first)")
	}

	// Тестируем пагинацию
	limitedPosts, err := storage.GetPosts(ctx, 2, 0)
	if err != nil {
		t.Fatalf("Failed to get limited posts: %v", err)
	}
	if len(limitedPosts) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(limitedPosts))
	}

	// Тестируем offset
	offsetPosts, err := storage.GetPosts(ctx, 2, 1)
	if err != nil {
		t.Fatalf("Failed to get offset posts: %v", err)
	}
	if len(offsetPosts) != 2 {
		t.Errorf("Expected 2 posts with offset, got %d", len(offsetPosts))
	}
}

// TestMemoryStorage_UpdatePost тестирует обновление поста
func TestMemoryStorage_UpdatePost(t *testing.T) {
	storage := repository.NewMemoryStorage()
	defer storage.Close()

	ctx := context.Background()

	// Создаем пост
	post := &model.Post{
		Title:   "Original Title",
		Content: "Original Content",
	}
	createdPost, err := storage.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}

	originalTime := createdPost.CreatedAt

	// Обновляем пост
	updatedPost := &model.Post{
		ID:              createdPost.ID,
		Title:           "Updated Title",
		Content:         "Updated Content",
		CommentsEnabled: false,
	}

	result, err := storage.UpdatePost(ctx, updatedPost)
	if err != nil {
		t.Fatalf("Failed to update post: %v", err)
	}

	// Проверяем обновления
	if result.Title != "Updated Title" {
		t.Errorf("Expected title to be updated to 'Updated Title', got %q", result.Title)
	}
	if result.Content != "Updated Content" {
		t.Errorf("Expected content to be updated to 'Updated Content', got %q", result.Content)
	}
	if result.CommentsEnabled {
		t.Error("Expected comments to be disabled")
	}
	if !result.CreatedAt.Equal(originalTime) {
		t.Error("Expected created_at to remain unchanged")
	}
}

// TestMemoryStorage_DeletePost тестирует удаление поста и каскадное удаление комментариев
func TestMemoryStorage_DeletePost(t *testing.T) {
	storage := repository.NewMemoryStorage()
	defer storage.Close()

	ctx := context.Background()

	// Создаем пост
	post := &model.Post{
		Title:   "Post to Delete",
		Content: "This post will be deleted",
	}
	createdPost, err := storage.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}

	// Создаем комментарий к посту
	comment := &model.Comment{
		PostID:  createdPost.ID,
		Content: "Comment on post to be deleted",
	}
	_, err = storage.CreateComment(ctx, comment)
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}

	// Удаляем пост
	err = storage.DeletePost(ctx, createdPost.ID)
	if err != nil {
		t.Fatalf("Failed to delete post: %v", err)
	}

	// Проверяем, что пост удален
	_, err = storage.GetPost(ctx, createdPost.ID)
	if err != repository.ErrNotFound {
		t.Error("Expected post to be deleted")
	}

	// Проверяем, что комментарии тоже удалены (каскадное удаление)
	_, err = storage.GetCommentsByPostID(ctx, createdPost.ID)
	if err == nil {
		t.Error("Expected error when getting comments for deleted post")
	}
}

// TestMemoryStorage_TogglePostComments тестирует переключение комментариев для поста
func TestMemoryStorage_TogglePostComments(t *testing.T) {
	storage := repository.NewMemoryStorage()
	defer storage.Close()

	ctx := context.Background()

	// Создаем пост
	post := &model.Post{
		Title:   "Post for Comments Toggle",
		Content: "Testing comments toggle",
	}
	createdPost, err := storage.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}

	// По умолчанию комментарии включены
	if !createdPost.CommentsEnabled {
		t.Error("Expected comments to be enabled by default")
	}

	// Отключаем комментарии
	err = storage.TogglePostComments(ctx, createdPost.ID, false)
	if err != nil {
		t.Fatalf("Failed to toggle comments: %v", err)
	}

	// Проверяем, что комментарии отключены
	fetchedPost, err := storage.GetPost(ctx, createdPost.ID)
	if err != nil {
		t.Fatalf("Failed to get post: %v", err)
	}
	if fetchedPost.CommentsEnabled {
		t.Error("Expected comments to be disabled")
	}

	// Включаем комментарии обратно
	err = storage.TogglePostComments(ctx, createdPost.ID, true)
	if err != nil {
		t.Fatalf("Failed to toggle comments back: %v", err)
	}

	// Проверяем, что комментарии включены
	fetchedPost, err = storage.GetPost(ctx, createdPost.ID)
	if err != nil {
		t.Fatalf("Failed to get post: %v", err)
	}
	if !fetchedPost.CommentsEnabled {
		t.Error("Expected comments to be enabled")
	}
}

// TestMemoryStorage_CreateComment тестирует создание комментария
func TestMemoryStorage_CreateComment(t *testing.T) {
	storage := repository.NewMemoryStorage()
	defer storage.Close()

	ctx := context.Background()

	// Создаем пост
	post := &model.Post{
		Title:   "Post for Comment",
		Content: "Post that will have comments",
	}
	createdPost, err := storage.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}

	// Создаем комментарий
	comment := &model.Comment{
		PostID:  createdPost.ID,
		Content: "This is a test comment",
	}
	createdComment, err := storage.CreateComment(ctx, comment)
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}

	// Проверяем комментарий
	if createdComment.ID == uuid.Nil {
		t.Error("Expected comment to have ID")
	}
	if createdComment.PostID != createdPost.ID {
		t.Error("Expected comment to be linked to the post")
	}
	if createdComment.Content != comment.Content {
		t.Errorf("Expected content %q, got %q", comment.Content, createdComment.Content)
	}
	if createdComment.ParentID != nil {
		t.Error("Expected root comment to have nil parent")
	}
	if createdComment.CreatedAt.IsZero() {
		t.Error("Expected created_at to be set")
	}
}

// TestMemoryStorage_CommentHierarchy тестирует иерархическую структуру комментариев
func TestMemoryStorage_CommentHierarchy(t *testing.T) {
	storage := repository.NewMemoryStorage()
	defer storage.Close()

	ctx := context.Background()

	// Создаем пост
	post := &model.Post{
		Title:   "Post with Comment Hierarchy",
		Content: "This post will have nested comments",
	}
	createdPost, err := storage.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}

	// Создаем корневой комментарий
	rootComment := &model.Comment{
		PostID:  createdPost.ID,
		Content: "Root comment",
	}
	createdRootComment, err := storage.CreateComment(ctx, rootComment)
	if err != nil {
		t.Fatalf("Failed to create root comment: %v", err)
	}

	// Создаем дочерний комментарий
	childComment := &model.Comment{
		PostID:   createdPost.ID,
		ParentID: &createdRootComment.ID,
		Content:  "Child comment",
	}
	createdChildComment, err := storage.CreateComment(ctx, childComment)
	if err != nil {
		t.Fatalf("Failed to create child comment: %v", err)
	}

	// Создаем внучатый комментарий
	grandChildComment := &model.Comment{
		PostID:   createdPost.ID,
		ParentID: &createdChildComment.ID,
		Content:  "Grandchild comment",
	}
	_, err = storage.CreateComment(ctx, grandChildComment)
	if err != nil {
		t.Fatalf("Failed to create grandchild comment: %v", err)
	}

	// Получаем все комментарии (плоский список)
	comments, err := storage.GetCommentsByPostID(ctx, createdPost.ID)
	if err != nil {
		t.Fatalf("Failed to get comments: %v", err)
	}
	if len(comments) != 3 {
		t.Errorf("Expected 3 comments, got %d", len(comments))
	}

	// Получаем иерархическое дерево
	tree, err := storage.GetCommentTree(ctx, createdPost.ID)
	if err != nil {
		t.Fatalf("Failed to get comment tree: %v", err)
	}

	// Проверяем структуру дерева
	if len(tree) != 1 {
		t.Errorf("Expected 1 root comment, got %d", len(tree))
	}

	rootNode := tree[0]
	if rootNode.Comment.Content != "Root comment" {
		t.Error("Expected root comment content to match")
	}
	if len(rootNode.Children) != 1 {
		t.Errorf("Expected 1 child comment, got %d", len(rootNode.Children))
	}

	childNode := rootNode.Children[0]
	if childNode.Comment.Content != "Child comment" {
		t.Error("Expected child comment content to match")
	}
	if len(childNode.Children) != 1 {
		t.Errorf("Expected 1 grandchild comment, got %d", len(childNode.Children))
	}

	grandChildNode := childNode.Children[0]
	if grandChildNode.Comment.Content != "Grandchild comment" {
		t.Error("Expected grandchild comment content to match")
	}
	if len(grandChildNode.Children) != 0 {
		t.Error("Expected grandchild to have no children")
	}
}

// TestMemoryStorage_CommentValidation тестирует валидацию комментариев
func TestMemoryStorage_CommentValidation(t *testing.T) {
	storage := repository.NewMemoryStorage()
	defer storage.Close()

	ctx := context.Background()

	// Создаем пост
	post := &model.Post{
		Title:   "Post for Validation Tests",
		Content: "Testing comment validation",
	}
	createdPost, err := storage.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}

	// Тест: комментарий к несуществующему посту
	invalidComment := &model.Comment{
		PostID:  uuid.New(), // Несуществующий пост
		Content: "Comment to non-existent post",
	}
	_, err = storage.CreateComment(ctx, invalidComment)
	if err == nil {
		t.Error("Expected error when creating comment for non-existent post")
	}

	// Отключаем комментарии для поста
	err = storage.TogglePostComments(ctx, createdPost.ID, false)
	if err != nil {
		t.Fatalf("Failed to disable comments: %v", err)
	}

	// Тест: комментарий к посту с отключенными комментариями
	disabledComment := &model.Comment{
		PostID:  createdPost.ID,
		Content: "Comment to post with disabled comments",
	}
	_, err = storage.CreateComment(ctx, disabledComment)
	if err == nil {
		t.Error("Expected error when creating comment for post with disabled comments")
	}

	// Включаем комментарии обратно
	err = storage.TogglePostComments(ctx, createdPost.ID, true)
	if err != nil {
		t.Fatalf("Failed to enable comments: %v", err)
	}

	// Тест: комментарий с несуществующим родителем
	orphanComment := &model.Comment{
		PostID:   createdPost.ID,
		ParentID: func() *uuid.UUID { id := uuid.New(); return &id }(), // Несуществующий родитель
		Content:  "Orphan comment",
	}
	_, err = storage.CreateComment(ctx, orphanComment)
	if err == nil {
		t.Error("Expected error when creating comment with non-existent parent")
	}
}

// TestMemoryStorage_DeleteComment тестирует удаление комментария с каскадным удалением
func TestMemoryStorage_DeleteComment(t *testing.T) {
	storage := repository.NewMemoryStorage()
	defer storage.Close()

	ctx := context.Background()

	// Создаем пост
	post := &model.Post{
		Title:   "Post for Comment Deletion",
		Content: "Testing comment deletion",
	}
	createdPost, err := storage.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}

	// Создаем иерархию комментариев
	rootComment, err := storage.CreateComment(ctx, &model.Comment{
		PostID:  createdPost.ID,
		Content: "Root comment to delete",
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

	// Проверяем, что у нас 3 комментария
	initialComments, err := storage.GetCommentsByPostID(ctx, createdPost.ID)
	if err != nil {
		t.Fatalf("Failed to get initial comments: %v", err)
	}
	if len(initialComments) != 3 {
		t.Errorf("Expected 3 comments initially, got %d", len(initialComments))
	}

	// Удаляем корневой комментарий (должен удалить все дочерние)
	err = storage.DeleteComment(ctx, rootComment.ID)
	if err != nil {
		t.Fatalf("Failed to delete root comment: %v", err)
	}

	// Проверяем, что все комментарии удалены
	remainingComments, err := storage.GetCommentsByPostID(ctx, createdPost.ID)
	if err != nil {
		t.Fatalf("Failed to get remaining comments: %v", err)
	}
	if len(remainingComments) != 0 {
		t.Errorf("Expected 0 comments after deletion, got %d", len(remainingComments))
	}
}

// TestMemoryStorage_HealthCheck тестирует health check
func TestMemoryStorage_HealthCheck(t *testing.T) {
	storage := repository.NewMemoryStorage()
	ctx := context.Background()

	// Тест: здоровое хранилище
	err := storage.HealthCheck(ctx)
	if err != nil {
		t.Errorf("Expected health check to pass, got error: %v", err)
	}

	// Закрываем хранилище
	storage.Close()

	// Тест: закрытое хранилище
	err = storage.HealthCheck(ctx)
	if err != repository.ErrConnectionFailed {
		t.Errorf("Expected ErrConnectionFailed after close, got: %v", err)
	}
}

// TestMemoryStorage_ConcurrentAccess тестирует concurrent доступ
func TestMemoryStorage_ConcurrentAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent access test in short mode")
	}

	storage := repository.NewMemoryStorage()
	defer storage.Close()

	ctx := context.Background()

	// Создаем пост для тестирования
	post := &model.Post{
		Title:   "Concurrent Test Post",
		Content: "Testing concurrent access",
	}
	createdPost, err := storage.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}

	// Запускаем несколько горутин для создания комментариев
	const numGoroutines = 10
	const commentsPerGoroutine = 5

	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			for j := 0; j < commentsPerGoroutine; j++ {
				comment := &model.Comment{
					PostID:  createdPost.ID,
					Content: fmt.Sprintf("Comment %d-%d", goroutineID, j),
				}
				_, err := storage.CreateComment(ctx, comment)
				if err != nil {
					results <- err
					return
				}
			}
			results <- nil
		}(i)
	}

	// Ждем завершения всех горутин
	for i := 0; i < numGoroutines; i++ {
		err := <-results
		if err != nil {
			t.Errorf("Goroutine error: %v", err)
		}
	}

	// Проверяем, что все комментарии созданы
	comments, err := storage.GetCommentsByPostID(ctx, createdPost.ID)
	if err != nil {
		t.Fatalf("Failed to get comments: %v", err)
	}

	expectedCount := numGoroutines * commentsPerGoroutine
	if len(comments) != expectedCount {
		t.Errorf("Expected %d comments, got %d", expectedCount, len(comments))
	}
}

// Benchmark тесты

// BenchmarkMemoryStorage_CreatePost бенчмарк создания постов
func BenchmarkMemoryStorage_CreatePost(b *testing.B) {
	storage := repository.NewMemoryStorage()
	defer storage.Close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		post := &model.Post{
			Title:   fmt.Sprintf("Benchmark Post %d", i),
			Content: fmt.Sprintf("Benchmark content %d", i),
		}
		_, err := storage.CreatePost(ctx, post)
		if err != nil {
			b.Fatalf("Failed to create post: %v", err)
		}
	}
}

// BenchmarkMemoryStorage_GetPosts бенчмарк получения постов
func BenchmarkMemoryStorage_GetPosts(b *testing.B) {
	storage := repository.NewMemoryStorage()
	defer storage.Close()

	ctx := context.Background()

	// Предварительно создаем посты
	for i := 0; i < 1000; i++ {
		post := &model.Post{
			Title:   fmt.Sprintf("Benchmark Post %d", i),
			Content: fmt.Sprintf("Benchmark content %d", i),
		}
		_, err := storage.CreatePost(ctx, post)
		if err != nil {
			b.Fatalf("Failed to create post: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := storage.GetPosts(ctx, 10, 0)
		if err != nil {
			b.Fatalf("Failed to get posts: %v", err)
		}
	}
}
