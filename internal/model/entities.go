// Package model содержит доменные модели (entities) системы комментариев.
// Эти модели представляют основные бизнес-сущности и не зависят от внешних слоев.
//
// Основные принципы:
// - Чистые бизнес-модели без зависимостей от базы данных или фреймворков
// - Доменная логика и валидация находятся в методах моделей
// - Все поля имеют осмысленные типы и ограничения
// - Поддержка JSON сериализации для API
package model

import (
	"time"

	"github.com/google/uuid"
)

// Post представляет пост в системе комментариев.
// Содержит основную информацию о посте и настройки комментирования.
//
// Бизнес-правила:
//   - Заголовок: обязательный, 1-255 символов
//   - Содержимое: обязательное, 1-10000 символов
//   - CommentsEnabled: управляет возможностью добавления комментариев
//   - ID генерируется автоматически при создании
//   - CreatedAt устанавливается в UTC при создании
//
// Валидация:
//
//	Используйте метод IsValid() для проверки корректности данных
//	перед сохранением в репозиторий.
//
// Пример использования:
//
//	post := &Post{
//		Title:   "Заголовок поста",
//		Content: "Содержимое поста",
//	}
//	if post.IsValid() {
//		post.Prepare() // Устанавливает ID и время
//		// Сохранение через репозиторий
//	}
type Post struct {
	ID              uuid.UUID `json:"id" db:"id"`                            // Уникальный идентификатор поста
	Title           string    `json:"title" db:"title"`                      // Заголовок поста (1-255 символов)
	Content         string    `json:"content" db:"content"`                  // Содержимое поста (до 10000 символов)
	CommentsEnabled bool      `json:"commentsEnabled" db:"comments_enabled"` // Флаг разрешения комментирования
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`             // Время создания поста (UTC)
}

// Comment представляет комментарий к посту.
// Поддерживает иерархическую структуру через ParentID.
//
// Бизнес-правила:
//   - Содержимое: обязательное, 1-2000 символов
//   - PostID: обязательный, ссылка на существующий пост
//   - ParentID: опциональный, для создания иерархии комментариев
//   - ID генерируется автоматически
//   - CreatedAt устанавливается в UTC при создании
//
// Иерархия:
//   - ParentID == nil: корневой комментарий
//   - ParentID != nil: ответ на комментарий с указанным ID
type Comment struct {
	ID        uuid.UUID  `json:"id" db:"id"`                        // Уникальный идентификатор комментария
	PostID    uuid.UUID  `json:"postId" db:"post_id"`               // ID поста, к которому относится комментарий
	ParentID  *uuid.UUID `json:"parentId,omitempty" db:"parent_id"` // ID родительского комментария (NULL для корневых)
	Content   string     `json:"content" db:"content"`              // Текст комментария (до 2000 символов)
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`         // Время создания комментария (UTC)
}

// PostWithComments объединяет пост с его комментариями.
// Используется для передачи полной информации о посте с комментариями.
type PostWithComments struct {
	Post               // Встроенная структура Post
	Comments []Comment `json:"comments"` // Список всех комментариев поста (плоский)
}

// CommentTree представляет иерархическую структуру комментариев.
// Позволяет построить дерево ответов для отображения вложенных комментариев.
type CommentTree struct {
	Comment                // Встроенная структура Comment
	Children []CommentTree `json:"children,omitempty"` // Дочерние комментарии (рекурсивная структура)
}

// Доменные методы для Post

// IsValidTitle проверяет валидность заголовка поста.
// Заголовок должен быть от 1 до 255 символов.
func (p *Post) IsValidTitle() bool {
	return len(p.Title) > 0 && len(p.Title) <= 255
}

// IsValidContent проверяет валидность содержимого поста.
// Содержимое не должно быть пустым и не должно превышать 10000 символов.
func (p *Post) IsValidContent() bool {
	return len(p.Content) > 0 && len(p.Content) <= 10000
}

// CanAddComments проверяет, можно ли добавлять комментарии к посту.
func (p *Post) CanAddComments() bool {
	return p.CommentsEnabled
}

// IsValid выполняет полную валидацию поста.
// Проверяет все обязательные поля и бизнес-правила.
func (p *Post) IsValid() bool {
	return p.IsValidTitle() && p.IsValidContent()
}

// Prepare подготавливает пост к сохранению.
// Устанавливает ID и время создания, если они не заданы.
func (p *Post) Prepare() {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
}

// Доменные методы для Comment

// IsValidComment проверяет валидность комментария.
// Содержимое должно быть от 1 до 2000 символов.
func (c *Comment) IsValidComment() bool {
	return len(c.Content) > 0 && len(c.Content) <= 2000
}

// IsRootComment проверяет, является ли комментарий корневым.
func (c *Comment) IsRootComment() bool {
	return c.ParentID == nil
}

// HasValidPost проверяет, что комментарий ссылается на валидный пост.
func (c *Comment) HasValidPost() bool {
	return c.PostID != uuid.Nil
}

// IsValid выполняет полную валидацию комментария.
func (c *Comment) IsValid() bool {
	return c.IsValidComment() && c.HasValidPost()
}

// Prepare подготавливает комментарий к сохранению.
// Устанавливает ID и время создания, если они не заданы.
func (c *Comment) Prepare() {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now().UTC()
	}
}

// Доменные методы для CommentTree

// HasChildren проверяет, есть ли у комментария дочерние комментарии.
func (ct *CommentTree) HasChildren() bool {
	return len(ct.Children) > 0
}

// GetChildrenCount возвращает количество прямых дочерних комментариев.
func (ct *CommentTree) GetChildrenCount() int {
	return len(ct.Children)
}
