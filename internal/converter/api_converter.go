// Package converter предоставляет конвертеры для преобразования данных между слоями.
// Отвечает за изоляцию между API слоем и доменными моделями.
//
// Основные компоненты:
// - GraphQLConverter: преобразование между GraphQL и доменными моделями
// - ValidationConverter: валидация входных данных API
//
// Принципы:
// - Четкое разделение ответственности
// - Централизованная валидация
// - Использование конфигурационных констант
// - Подробные сообщения об ошибках
package converter

import (
	"fmt"
	"time"

	"github.com/NarthurN/CommentsSystem/internal/config"
	"github.com/NarthurN/CommentsSystem/internal/model"
	"github.com/google/uuid"
)

// GraphQLConverter отвечает за преобразование данных между GraphQL представлением и доменными моделями.
// Инкапсулирует логику конвертации и форматирования для GraphQL API.
//
// Основные возможности:
// - Преобразование Post в GraphQL формат
// - Преобразование Comment в GraphQL формат
// - Форматирование дат в ISO 8601
// - Обработка nullable полей
type GraphQLConverter struct {
	config *config.Config // Конфигурация для лимитов и форматирования
}

// NewGraphQLConverter создает новый экземпляр GraphQL конвертера.
// Принимает конфигурацию для использования настроек лимитов и форматирования.
func NewGraphQLConverter(cfg *config.Config) *GraphQLConverter {
	return &GraphQLConverter{
		config: cfg,
	}
}

// PostToGraphQL преобразует доменную модель Post в формат для GraphQL.
// Выполняет форматирование и валидацию данных.
//
// Параметры:
//   - post: доменная модель поста
//
// Возвращает:
//   - map[string]interface{}: пост в формате GraphQL
//   - error: ошибка преобразования
func (c *GraphQLConverter) PostToGraphQL(post *model.Post) (map[string]interface{}, error) {
	if post == nil {
		return nil, ErrNilPost
	}

	// Валидируем пост перед конвертацией
	if !post.IsValid() {
		return nil, ErrInvalidPost
	}

	return map[string]interface{}{
		"id":              post.ID.String(),
		"title":           post.Title,
		"content":         post.Content,
		"commentsEnabled": post.CommentsEnabled,
		"createdAt":       post.CreatedAt.Format(time.RFC3339),
	}, nil
}

// CommentToGraphQL преобразует доменную модель Comment в формат для GraphQL.
// Обрабатывает nullable поля и форматирование.
//
// Параметры:
//   - comment: доменная модель комментария
//
// Возвращает:
//   - map[string]interface{}: комментарий в формате GraphQL
//   - error: ошибка преобразования
func (c *GraphQLConverter) CommentToGraphQL(comment *model.Comment) (map[string]interface{}, error) {
	if comment == nil {
		return nil, ErrNilComment
	}

	// Валидируем комментарий перед конвертацией
	if !comment.IsValid() {
		return nil, ErrInvalidComment
	}

	result := map[string]interface{}{
		"id":        comment.ID.String(),
		"postId":    comment.PostID.String(),
		"content":   comment.Content,
		"createdAt": comment.CreatedAt.Format(time.RFC3339),
	}

	// Обрабатываем nullable поле ParentID
	if comment.ParentID != nil {
		result["parentId"] = comment.ParentID.String()
	} else {
		result["parentId"] = nil
	}

	return result, nil
}

// PostsToGraphQL преобразует срез постов в формат для GraphQL.
// Обрабатывает пагинацию и валидацию.
//
// Параметры:
//   - posts: срез доменных моделей постов
//
// Возвращает:
//   - []map[string]interface{}: посты в формате GraphQL
//   - error: ошибка преобразования
func (c *GraphQLConverter) PostsToGraphQL(posts []*model.Post) ([]map[string]interface{}, error) {
	if posts == nil {
		return []map[string]interface{}{}, nil
	}

	result := make([]map[string]interface{}, 0, len(posts))
	for _, post := range posts {
		graphqlPost, err := c.PostToGraphQL(post)
		if err != nil {
			return nil, fmt.Errorf("failed to convert post %s: %w", post.ID, err)
		}
		result = append(result, graphqlPost)
	}

	return result, nil
}

// CommentsToGraphQL преобразует срез комментариев в формат для GraphQL.
// Обрабатывает пагинацию и валидацию.
//
// Параметры:
//   - comments: срез доменных моделей комментариев
//
// Возвращает:
//   - []map[string]interface{}: комментарии в формате GraphQL
//   - error: ошибка преобразования
func (c *GraphQLConverter) CommentsToGraphQL(comments []model.Comment) ([]map[string]interface{}, error) {
	if comments == nil {
		return []map[string]interface{}{}, nil
	}

	result := make([]map[string]interface{}, 0, len(comments))
	for _, comment := range comments {
		graphqlComment, err := c.CommentToGraphQL(&comment)
		if err != nil {
			return nil, fmt.Errorf("failed to convert comment %s: %w", comment.ID, err)
		}
		result = append(result, graphqlComment)
	}

	return result, nil
}

// ValidationConverter отвечает за валидацию входных данных от API.
// Использует конфигурационные константы для проверки лимитов.
//
// Основные возможности:
// - Валидация создания постов
// - Валидация создания комментариев
// - Проверка UUID форматов
// - Использование настраиваемых лимитов
type ValidationConverter struct {
	config *config.Config // Конфигурация с лимитами валидации
}

// NewValidationConverter создает новый экземпляр валидационного конвертера.
// Принимает конфигурацию для использования настраиваемых лимитов.
func NewValidationConverter(cfg *config.Config) *ValidationConverter {
	return &ValidationConverter{
		config: cfg,
	}
}

// ValidateAndConvertCreatePost валидирует входные данные для создания поста
// и конвертирует их в доменную модель.
//
// Параметры:
//   - title: заголовок поста
//   - content: содержимое поста
//
// Возвращает:
//   - *model.Post: готовая к сохранению доменная модель
//   - error: ошибка валидации или конвертации
func (c *ValidationConverter) ValidateAndConvertCreatePost(title, content string) (*model.Post, error) {
	// Валидируем входные данные
	if err := c.validatePostInput(title, content); err != nil {
		return nil, err
	}

	// Создаем доменную модель
	post := &model.Post{
		Title:           title,
		Content:         content,
		CommentsEnabled: true, // По умолчанию комментирование включено
	}

	// Подготавливаем к сохранению (устанавливаем ID и время)
	post.Prepare()

	return post, nil
}

// ValidateAndConvertCreateComment валидирует входные данные для создания комментария
// и конвертирует их в доменную модель.
//
// Параметры:
//   - postID: ID поста для комментария
//   - parentID: ID родительского комментария (может быть пустым)
//   - content: текст комментария
//
// Возвращает:
//   - *model.Comment: готовая к сохранению доменная модель
//   - error: ошибка валидации или конвертации
func (c *ValidationConverter) ValidateAndConvertCreateComment(postID, parentID, content string) (*model.Comment, error) {
	// Валидируем базовые входные данные
	if err := c.validateCommentInput(postID, content); err != nil {
		return nil, err
	}

	// Парсим обязательный PostID
	postUUID, err := parseUUID(postID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPostID, err)
	}

	// Создаем доменную модель
	comment := &model.Comment{
		PostID:  postUUID,
		Content: content,
	}

	// Парсим опциональный ParentID
	if parentID != "" {
		parentUUID, err := parseUUID(parentID)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidParentID, err)
		}
		comment.ParentID = &parentUUID
	}

	// Подготавливаем к сохранению (устанавливаем ID и время)
	comment.Prepare()

	return comment, nil
}

// ValidatePaginationParams валидирует параметры пагинации.
// Проверяет и нормализует значения limit и offset.
//
// Параметры:
//   - limit: максимальное количество элементов (может быть nil)
//   - offset: количество пропускаемых элементов (может быть nil)
//   - defaultLimit: лимит по умолчанию
//
// Возвращает:
//   - int: валидное значение limit
//   - int: валидное значение offset
//   - error: ошибка валидации
func (c *ValidationConverter) ValidatePaginationParams(limit, offset *int, defaultLimit int) (int, int, error) {
	// Устанавливаем значения по умолчанию
	resultLimit := defaultLimit
	resultOffset := 0

	// Валидируем и устанавливаем limit
	if limit != nil {
		if *limit <= 0 {
			return 0, 0, ErrInvalidLimit
		}
		if *limit > 100 { // Максимальный лимит для защиты от злоупотреблений
			return 0, 0, ErrLimitTooLarge
		}
		resultLimit = *limit
	}

	// Валидируем и устанавливаем offset
	if offset != nil {
		if *offset < 0 {
			return 0, 0, ErrInvalidOffset
		}
		resultOffset = *offset
	}

	return resultLimit, resultOffset, nil
}

// validatePostInput выполняет валидацию входных данных поста.
// Проверяет длину заголовка и содержимого согласно конфигурации.
func (c *ValidationConverter) validatePostInput(title, content string) error {
	if title == "" {
		return ErrEmptyTitle
	}
	if len(title) > c.config.MaxTitleLength {
		return fmt.Errorf("%w: максимум %d символов, получено %d",
			ErrTitleTooLong, c.config.MaxTitleLength, len(title))
	}
	if content == "" {
		return ErrEmptyContent
	}
	if len(content) > c.config.MaxContentLength {
		return fmt.Errorf("%w: максимум %d символов, получено %d",
			ErrContentTooLong, c.config.MaxContentLength, len(content))
	}
	return nil
}

// validateCommentInput выполняет валидацию входных данных комментария.
// Проверяет длину содержимого согласно конфигурации.
func (c *ValidationConverter) validateCommentInput(postID, content string) error {
	if postID == "" {
		return ErrEmptyPostID
	}
	if content == "" {
		return ErrEmptyContent
	}
	if len(content) > c.config.MaxCommentLength {
		return fmt.Errorf("%w: максимум %d символов, получено %d",
			ErrContentTooLong, c.config.MaxCommentLength, len(content))
	}
	return nil
}

// parseUUID безопасно парсит UUID строку.
// Возвращает подробные ошибки для отладки.
func parseUUID(s string) (uuid.UUID, error) {
	if s == "" {
		return uuid.Nil, fmt.Errorf("пустая строка UUID")
	}

	parsed, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("некорректный формат UUID '%s': %w", s, err)
	}

	return parsed, nil
}

// Ошибки валидации и конвертации
var (
	// Ошибки данных
	ErrNilPost        = fmt.Errorf("пост не может быть nil")
	ErrNilComment     = fmt.Errorf("комментарий не может быть nil")
	ErrInvalidPost    = fmt.Errorf("пост содержит некорректные данные")
	ErrInvalidComment = fmt.Errorf("комментарий содержит некорректные данные")

	// Ошибки валидации полей
	ErrInvalidTitle    = fmt.Errorf("некорректный заголовок")
	ErrInvalidContent  = fmt.Errorf("некорректное содержимое")
	ErrInvalidPostID   = fmt.Errorf("некорректный ID поста")
	ErrInvalidParentID = fmt.Errorf("некорректный ID родительского комментария")

	// Ошибки пустых значений
	ErrEmptyTitle   = fmt.Errorf("заголовок не может быть пустым")
	ErrEmptyContent = fmt.Errorf("содержимое не может быть пустым")
	ErrEmptyPostID  = fmt.Errorf("ID поста не может быть пустым")

	// Ошибки лимитов
	ErrTitleTooLong   = fmt.Errorf("заголовок слишком длинный")
	ErrContentTooLong = fmt.Errorf("содержимое слишком длинное")

	// Ошибки пагинации
	ErrInvalidLimit  = fmt.Errorf("лимит должен быть положительным числом")
	ErrInvalidOffset = fmt.Errorf("смещение не может быть отрицательным")
	ErrLimitTooLarge = fmt.Errorf("лимит слишком большой")
)
