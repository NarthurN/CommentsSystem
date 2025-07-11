package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PostDB представляет модель поста в базе данных
// Отражает структуру таблицы posts
type PostDB struct {
	ID              uuid.UUID `db:"id"`
	Title           string    `db:"title"`
	Content         string    `db:"content"`
	CommentsEnabled bool      `db:"comments_enabled"`
	CreatedAt       time.Time `db:"created_at"`
}

// CommentDB представляет модель комментария в базе данных
// Отражает структуру таблицы comments
type CommentDB struct {
	ID        uuid.UUID  `db:"id"`
	PostID    uuid.UUID  `db:"post_id"`
	ParentID  *uuid.UUID `db:"parent_id"`
	Content   string     `db:"content"`
	CreatedAt time.Time  `db:"created_at"`
}

// PostWithCommentsDB представляет пост с комментариями для JOIN запросов
type PostWithCommentsDB struct {
	PostDB
	// Дополнительные поля для JOIN
	CommentID        *uuid.UUID `db:"comment_id"`
	CommentPostID    *uuid.UUID `db:"comment_post_id"`
	CommentParentID  *uuid.UUID `db:"comment_parent_id"`
	CommentContent   *string    `db:"comment_content"`
	CommentCreatedAt *time.Time `db:"comment_created_at"`
}

// CommentTreeDB представляет результат рекурсивного CTE запроса
type CommentTreeDB struct {
	CommentDB
	Level int `db:"level"` // Уровень вложенности из CTE
}

// TableName возвращает имя таблицы для PostDB
func (PostDB) TableName() string {
	return "posts"
}

// TableName возвращает имя таблицы для CommentDB
func (CommentDB) TableName() string {
	return "comments"
}

// Методы для работы с базой данных

// GetSelectColumns возвращает список колонок для SELECT запроса постов
func (PostDB) GetSelectColumns() []string {
	return []string{"id", "title", "content", "comments_enabled", "created_at"}
}

// GetSelectColumns возвращает список колонок для SELECT запроса комментариев
func (CommentDB) GetSelectColumns() []string {
	return []string{"id", "post_id", "parent_id", "content", "created_at"}
}

// GetInsertColumns возвращает список колонок для INSERT запроса постов
func (PostDB) GetInsertColumns() []string {
	return []string{"id", "title", "content", "comments_enabled", "created_at"}
}

// GetInsertColumns возвращает список колонок для INSERT запроса комментариев
func (CommentDB) GetInsertColumns() []string {
	return []string{"id", "post_id", "parent_id", "content", "created_at"}
}

// GetUpdateColumns возвращает список колонок для UPDATE запроса постов
func (PostDB) GetUpdateColumns() []string {
	return []string{"title", "content", "comments_enabled"}
}

// Validation методы на уровне репозитория

// Validate проверяет валидность модели поста на уровне БД
func (p *PostDB) Validate() error {
	if p.ID == uuid.Nil {
		return ErrInvalidID
	}
	if len(p.Title) == 0 || len(p.Title) > 255 {
		return ErrInvalidTitle
	}
	if len(p.Content) == 0 {
		return ErrInvalidContent
	}
	return nil
}

// Validate проверяет валидность модели комментария на уровне БД
func (c *CommentDB) Validate() error {
	if c.ID == uuid.Nil {
		return ErrInvalidID
	}
	if c.PostID == uuid.Nil {
		return ErrInvalidPostID
	}
	if len(c.Content) == 0 || len(c.Content) > 2000 {
		return ErrInvalidContent
	}
	return nil
}

// Ошибки валидации репозитория
var (
	ErrInvalidID      = fmt.Errorf("invalid ID")
	ErrInvalidTitle   = fmt.Errorf("invalid title")
	ErrInvalidContent = fmt.Errorf("invalid content")
	ErrInvalidPostID  = fmt.Errorf("invalid post ID")
)
