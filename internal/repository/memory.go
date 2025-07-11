package repository

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/NarthurN/CommentsSystem/internal/model"
	"github.com/google/uuid"
)

// MemoryStorage реализует интерфейс Storage для in-memory хранилища.
// Обеспечивает thread-safe операции с данными в памяти.
//
// Особенности:
// - Thread-safe операции через sync.RWMutex
// - Автоматическая генерация ID и времени создания
// - Поддержка всех операций интерфейса Storage
// - Имитация поведения реальной базы данных
// - Сортировка постов по времени создания (новые первыми)
// - Каскадное удаление комментариев при удалении поста
type MemoryStorage struct {
	mu       sync.RWMutex                 // Мьютекс для thread-safe операций
	posts    map[uuid.UUID]*model.Post    // Хранилище постов
	comments map[uuid.UUID]*model.Comment // Хранилище комментариев
	closed   bool                         // Флаг закрытия хранилища
}

// NewMemoryStorage создает новый экземпляр in-memory хранилища.
// Инициализирует внутренние структуры данных и готов к использованию.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		posts:    make(map[uuid.UUID]*model.Post),
		comments: make(map[uuid.UUID]*model.Comment),
		closed:   false,
	}
}

// Close закрывает хранилище и очищает данные.
// После вызова Close хранилище становится недоступным для операций.
func (s *MemoryStorage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	// Очищаем данные
	s.posts = nil
	s.comments = nil
	s.closed = true

	return nil
}

// HealthCheck проверяет состояние хранилища.
// Для in-memory хранилища всегда возвращает успех, если не закрыто.
func (s *MemoryStorage) HealthCheck(ctx context.Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return ErrConnectionFailed
	}

	return nil
}

// checkClosed проверяет, что хранилище не закрыто.
// Должно вызываться под мьютексом.
func (s *MemoryStorage) checkClosed() error {
	if s.closed {
		return ErrConnectionFailed
	}
	return nil
}

// Операции с постами

// CreatePost создает новый пост в памяти.
// Автоматически генерирует ID и устанавливает время создания.
func (s *MemoryStorage) CreatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkClosed(); err != nil {
		return nil, err
	}

	// Валидируем входные данные
	if post == nil {
		return nil, ErrInvalidInput
	}

	// Создаем копию для безопасности
	newPost := &model.Post{
		Title:           post.Title,
		Content:         post.Content,
		CommentsEnabled: true, // По умолчанию комментарии включены
	}

	// Если в переданном посте уже установлен флаг, используем его
	if post.ID != uuid.Nil && !post.CreatedAt.IsZero() {
		newPost.CommentsEnabled = post.CommentsEnabled
	}

	// Генерируем ID если не задан
	if newPost.ID == uuid.Nil {
		newPost.ID = uuid.New()
	}

	// Устанавливаем время создания если не задано
	if newPost.CreatedAt.IsZero() {
		newPost.CreatedAt = time.Now().UTC()
	}

	// Валидируем бизнес-правила
	if !newPost.IsValid() {
		return nil, ErrInvalidInput
	}

	// Проверяем уникальность ID
	if _, exists := s.posts[newPost.ID]; exists {
		return nil, ErrDuplicate
	}

	// Сохраняем пост
	s.posts[newPost.ID] = newPost

	// Возвращаем копию
	return &model.Post{
		ID:              newPost.ID,
		Title:           newPost.Title,
		Content:         newPost.Content,
		CommentsEnabled: newPost.CommentsEnabled,
		CreatedAt:       newPost.CreatedAt,
	}, nil
}

// GetPost получает пост по ID из памяти.
func (s *MemoryStorage) GetPost(ctx context.Context, id uuid.UUID) (*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := s.checkClosed(); err != nil {
		return nil, err
	}

	post, exists := s.posts[id]
	if !exists {
		return nil, ErrNotFound
	}

	// Возвращаем копию для безопасности
	return &model.Post{
		ID:              post.ID,
		Title:           post.Title,
		Content:         post.Content,
		CommentsEnabled: post.CommentsEnabled,
		CreatedAt:       post.CreatedAt,
	}, nil
}

// GetPosts получает список постов с пагинацией.
// Посты сортируются по времени создания (новые первыми).
func (s *MemoryStorage) GetPosts(ctx context.Context, limit, offset int) ([]*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := s.checkClosed(); err != nil {
		return nil, err
	}

	// Валидируем параметры пагинации
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Собираем все посты в слайс
	posts := make([]*model.Post, 0, len(s.posts))
	for _, post := range s.posts {
		posts = append(posts, post)
	}

	// Сортируем по времени создания (новые первыми)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})

	// Применяем пагинацию
	start := offset
	if start >= len(posts) {
		return []*model.Post{}, nil
	}

	end := start + limit
	if end > len(posts) {
		end = len(posts)
	}

	// Создаем копии для возврата
	result := make([]*model.Post, 0, end-start)
	for i := start; i < end; i++ {
		post := posts[i]
		result = append(result, &model.Post{
			ID:              post.ID,
			Title:           post.Title,
			Content:         post.Content,
			CommentsEnabled: post.CommentsEnabled,
			CreatedAt:       post.CreatedAt,
		})
	}

	return result, nil
}

// UpdatePost обновляет существующий пост в памяти.
func (s *MemoryStorage) UpdatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkClosed(); err != nil {
		return nil, err
	}

	if post == nil {
		return nil, ErrInvalidInput
	}

	// Проверяем, что пост существует
	existing, exists := s.posts[post.ID]
	if !exists {
		return nil, ErrNotFound
	}

	// Валидируем бизнес-правила
	if !post.IsValid() {
		return nil, ErrInvalidInput
	}

	// Обновляем пост (сохраняем время создания)
	updatedPost := &model.Post{
		ID:              post.ID,
		Title:           post.Title,
		Content:         post.Content,
		CommentsEnabled: post.CommentsEnabled,
		CreatedAt:       existing.CreatedAt, // Сохраняем оригинальное время
	}

	s.posts[post.ID] = updatedPost

	// Возвращаем копию
	return &model.Post{
		ID:              updatedPost.ID,
		Title:           updatedPost.Title,
		Content:         updatedPost.Content,
		CommentsEnabled: updatedPost.CommentsEnabled,
		CreatedAt:       updatedPost.CreatedAt,
	}, nil
}

// DeletePost удаляет пост и все связанные комментарии.
func (s *MemoryStorage) DeletePost(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkClosed(); err != nil {
		return err
	}

	// Проверяем, что пост существует
	if _, exists := s.posts[id]; !exists {
		return ErrNotFound
	}

	// Удаляем пост
	delete(s.posts, id)

	// Удаляем все комментарии к посту (каскадное удаление)
	for commentID, comment := range s.comments {
		if comment.PostID == id {
			delete(s.comments, commentID)
		}
	}

	return nil
}

// TogglePostComments включает/выключает комментарии для поста.
func (s *MemoryStorage) TogglePostComments(ctx context.Context, id uuid.UUID, enabled bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkClosed(); err != nil {
		return err
	}

	// Проверяем, что пост существует
	post, exists := s.posts[id]
	if !exists {
		return ErrNotFound
	}

	// Обновляем флаг комментариев
	post.CommentsEnabled = enabled

	return nil
}

// Операции с комментариями

// CreateComment создает новый комментарий в памяти.
func (s *MemoryStorage) CreateComment(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkClosed(); err != nil {
		return nil, err
	}

	if comment == nil {
		return nil, ErrInvalidInput
	}

	// Проверяем, что пост существует
	post, exists := s.posts[comment.PostID]
	if !exists {
		return nil, fmt.Errorf("post not found")
	}

	// Проверяем, что комментарии разрешены
	if !post.CommentsEnabled {
		return nil, fmt.Errorf("comments are disabled for this post")
	}

	// Если указан родительский комментарий, проверяем его существование
	if comment.ParentID != nil {
		parentComment, exists := s.comments[*comment.ParentID]
		if !exists {
			return nil, fmt.Errorf("parent comment not found")
		}
		// Проверяем, что родительский комментарий относится к тому же посту
		if parentComment.PostID != comment.PostID {
			return nil, fmt.Errorf("parent comment belongs to different post")
		}
	}

	// Создаем копию комментария
	newComment := &model.Comment{
		PostID:   comment.PostID,
		ParentID: comment.ParentID,
		Content:  comment.Content,
	}

	// Генерируем ID если не задан
	if newComment.ID == uuid.Nil {
		newComment.ID = uuid.New()
	}

	// Устанавливаем время создания если не задано
	if newComment.CreatedAt.IsZero() {
		newComment.CreatedAt = time.Now().UTC()
	}

	// Валидируем бизнес-правила
	if !newComment.IsValid() {
		return nil, ErrInvalidInput
	}

	// Проверяем уникальность ID
	if _, exists := s.comments[newComment.ID]; exists {
		return nil, ErrDuplicate
	}

	// Сохраняем комментарий
	s.comments[newComment.ID] = newComment

	// Возвращаем копию
	return &model.Comment{
		ID:        newComment.ID,
		PostID:    newComment.PostID,
		ParentID:  newComment.ParentID,
		Content:   newComment.Content,
		CreatedAt: newComment.CreatedAt,
	}, nil
}

// GetComment получает комментарий по ID.
func (s *MemoryStorage) GetComment(ctx context.Context, id uuid.UUID) (*model.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := s.checkClosed(); err != nil {
		return nil, err
	}

	comment, exists := s.comments[id]
	if !exists {
		return nil, ErrNotFound
	}

	// Возвращаем копию
	return &model.Comment{
		ID:        comment.ID,
		PostID:    comment.PostID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}, nil
}

// GetCommentsByPostID получает все комментарии для поста.
// Возвращает комментарии отсортированные по времени создания.
func (s *MemoryStorage) GetCommentsByPostID(ctx context.Context, postID uuid.UUID) ([]model.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := s.checkClosed(); err != nil {
		return nil, err
	}

	// Проверяем, что пост существует
	if _, exists := s.posts[postID]; !exists {
		return nil, fmt.Errorf("post not found")
	}

	// Собираем комментарии для поста
	var comments []model.Comment
	for _, comment := range s.comments {
		if comment.PostID == postID {
			comments = append(comments, model.Comment{
				ID:        comment.ID,
				PostID:    comment.PostID,
				ParentID:  comment.ParentID,
				Content:   comment.Content,
				CreatedAt: comment.CreatedAt,
			})
		}
	}

	// Сортируем по времени создания
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.Before(comments[j].CreatedAt)
	})

	return comments, nil
}

// GetCommentsByParentID получает дочерние комментарии с пагинацией
// ПРОИЗВОДИТЕЛЬНОСТЬ: Решает N+1 проблему в GraphQL children резолвере
func (s *MemoryStorage) GetCommentsByParentID(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]model.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := s.checkClosed(); err != nil {
		return nil, err
	}

	// Валидируем параметры пагинации
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Собираем дочерние комментарии
	var children []model.Comment
	for _, comment := range s.comments {
		if comment.ParentID != nil && *comment.ParentID == parentID {
			children = append(children, model.Comment{
				ID:        comment.ID,
				PostID:    comment.PostID,
				ParentID:  comment.ParentID,
				Content:   comment.Content,
				CreatedAt: comment.CreatedAt,
			})
		}
	}

	// Сортируем по времени создания
	sort.Slice(children, func(i, j int) bool {
		return children[i].CreatedAt.Before(children[j].CreatedAt)
	})

	// Применяем пагинацию
	start := offset
	if start >= len(children) {
		return []model.Comment{}, nil
	}

	end := start + limit
	if end > len(children) {
		end = len(children)
	}

	return children[start:end], nil
}

// GetRootCommentsByPostID получает только корневые комментарии с пагинацией
// ПРОИЗВОДИТЕЛЬНОСТЬ: Избегаем загрузки всех комментариев сразу
func (s *MemoryStorage) GetRootCommentsByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]model.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := s.checkClosed(); err != nil {
		return nil, err
	}

	// Проверяем, что пост существует
	if _, exists := s.posts[postID]; !exists {
		return nil, fmt.Errorf("post not found")
	}

	// Валидируем параметры пагинации
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Собираем только корневые комментарии для поста
	var rootComments []model.Comment
	for _, comment := range s.comments {
		if comment.PostID == postID && comment.ParentID == nil {
			rootComments = append(rootComments, model.Comment{
				ID:        comment.ID,
				PostID:    comment.PostID,
				ParentID:  comment.ParentID,
				Content:   comment.Content,
				CreatedAt: comment.CreatedAt,
			})
		}
	}

	// Сортируем по времени создания
	sort.Slice(rootComments, func(i, j int) bool {
		return rootComments[i].CreatedAt.Before(rootComments[j].CreatedAt)
	})

	// Применяем пагинацию
	start := offset
	if start >= len(rootComments) {
		return []model.Comment{}, nil
	}

	end := start + limit
	if end > len(rootComments) {
		end = len(rootComments)
	}

	return rootComments[start:end], nil
}

// DeleteComment удаляет комментарий и все дочерние комментарии.
func (s *MemoryStorage) DeleteComment(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkClosed(); err != nil {
		return err
	}

	// Проверяем, что комментарий существует
	if _, exists := s.comments[id]; !exists {
		return ErrNotFound
	}

	// Рекурсивно удаляем комментарий и всех его потомков
	s.deleteCommentRecursive(id)

	return nil
}

// deleteCommentRecursive рекурсивно удаляет комментарий и всех его потомков.
// Должно вызываться под мьютексом.
func (s *MemoryStorage) deleteCommentRecursive(id uuid.UUID) {
	// Сначала удаляем всех детей
	for commentID, comment := range s.comments {
		if comment.ParentID != nil && *comment.ParentID == id {
			s.deleteCommentRecursive(commentID)
		}
	}

	// Затем удаляем сам комментарий
	delete(s.comments, id)
}

// GetPostWithComments получает пост с комментариями (заглушка для совместимости).
func (s *MemoryStorage) GetPostWithComments(ctx context.Context, postID uuid.UUID) (*model.PostWithComments, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := s.checkClosed(); err != nil {
		return nil, err
	}

	// Получаем пост
	post, exists := s.posts[postID]
	if !exists {
		return nil, ErrNotFound
	}

	// Получаем комментарии
	var comments []model.Comment
	for _, comment := range s.comments {
		if comment.PostID == postID {
			comments = append(comments, *comment)
		}
	}

	// Сортируем комментарии по времени создания
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.Before(comments[j].CreatedAt)
	})

	return &model.PostWithComments{
		Post:     *post,
		Comments: comments,
	}, nil
}

// GetCommentTree получает иерархическое дерево комментариев для поста.
// Строит полную иерархию с рекурсивной вложенностью комментариев.
func (s *MemoryStorage) GetCommentTree(ctx context.Context, postID uuid.UUID) ([]model.CommentTree, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := s.checkClosed(); err != nil {
		return nil, err
	}

	// Получаем все комментарии для поста
	comments, err := s.GetCommentsByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Создаем карту для быстрого поиска комментариев по ID
	commentMap := make(map[uuid.UUID]model.Comment)
	for _, comment := range comments {
		commentMap[comment.ID] = comment
	}

	// Строим иерархию рекурсивно
	return s.buildCommentTree(commentMap, nil), nil
}

// buildCommentTree рекурсивно строит дерево комментариев.
// parentID = nil для корневых комментариев
func (s *MemoryStorage) buildCommentTree(commentMap map[uuid.UUID]model.Comment, parentID *uuid.UUID) []model.CommentTree {
	var children []model.CommentTree

	for _, comment := range commentMap {
		// Проверяем, является ли комментарий дочерним для указанного родителя
		if (parentID == nil && comment.ParentID == nil) ||
			(parentID != nil && comment.ParentID != nil && *comment.ParentID == *parentID) {

			// Рекурсивно строим детей для этого комментария
			childNodes := s.buildCommentTree(commentMap, &comment.ID)

			children = append(children, model.CommentTree{
				Comment:  comment,
				Children: childNodes,
			})
		}
	}

	// Сортируем по времени создания
	sort.Slice(children, func(i, j int) bool {
		return children[i].Comment.CreatedAt.Before(children[j].Comment.CreatedAt)
	})

	return children
}

// GetCommentHierarchy получает иерархию комментариев (алиас для GetCommentTree).
func (s *MemoryStorage) GetCommentHierarchy(ctx context.Context, postID uuid.UUID) ([]model.CommentTree, error) {
	return s.GetCommentTree(ctx, postID)
}
