package repository

import (
	"context"
	"errors"

	"github.com/NarthurN/CommentsSystem/internal/model"
	"github.com/google/uuid"
)

// Common repository errors
var (
	// ErrNotFound indicates that the requested entity was not found
	ErrNotFound = errors.New("entity not found")

	// ErrDuplicate indicates that the entity already exists
	ErrDuplicate = errors.New("entity already exists")

	// ErrInvalidInput indicates that the provided input is invalid
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnsupportedStorageType indicates that the storage type is not supported
	ErrUnsupportedStorageType = errors.New("unsupported storage type")

	// ErrConnectionFailed indicates that database connection failed
	ErrConnectionFailed = errors.New("database connection failed")

	// ErrTransactionFailed indicates that database transaction failed
	ErrTransactionFailed = errors.New("database transaction failed")
)

// Storage представляет интерфейс для работы с хранилищем данных.
// Этот интерфейс определяется в сервисном слое и реализуется в репозиторном слое
// в соответствии с принципами Dependency Inversion.
//
// Storage предоставляет унифицированный доступ к данным о постах и комментариях,
// скрывая детали реализации хранения от бизнес-логики.
type Storage interface {
	// Post operations

	// CreatePost создает новый пост в хранилище.
	// Возвращает созданный пост с заполненным ID и временем создания.
	CreatePost(ctx context.Context, post *model.Post) (*model.Post, error)

	// GetPost получает пост по ID.
	// Возвращает ErrNotFound если пост не найден.
	GetPost(ctx context.Context, id uuid.UUID) (*model.Post, error)

	// GetPosts получает список постов с пагинацией.
	// limit - максимальное количество постов
	// offset - количество пропускаемых постов
	GetPosts(ctx context.Context, limit, offset int) ([]*model.Post, error)

	// UpdatePost обновляет существующий пост.
	// Возвращает ErrNotFound если пост не найден.
	UpdatePost(ctx context.Context, post *model.Post) (*model.Post, error)

	// DeletePost удаляет пост по ID.
	// Также удаляет все связанные комментарии.
	DeletePost(ctx context.Context, id uuid.UUID) error

	// TogglePostComments включает/выключает возможность комментирования поста.
	TogglePostComments(ctx context.Context, id uuid.UUID, enabled bool) error

	// Comment operations

	// CreateComment создает новый комментарий.
	// Возвращает созданный комментарий с заполненным ID и временем создания.
	CreateComment(ctx context.Context, comment *model.Comment) (*model.Comment, error)

	// GetComment получает комментарий по ID.
	// Возвращает ErrNotFound если комментарий не найден.
	GetComment(ctx context.Context, id uuid.UUID) (*model.Comment, error)

	// GetCommentsByPostID получает все комментарии для поста.
	// Возвращает плоский список комментариев, отсортированный по времени создания.
	GetCommentsByPostID(ctx context.Context, postID uuid.UUID) ([]model.Comment, error)

	// GetCommentTree получает иерархическое дерево комментариев для поста.
	// Возвращает структурированное дерево с вложенными комментариями.
	GetCommentTree(ctx context.Context, postID uuid.UUID) ([]model.CommentTree, error)

	// DeleteComment удаляет комментарий и все его дочерние комментарии.
	DeleteComment(ctx context.Context, id uuid.UUID) error

	// Complex operations

	// GetPostWithComments получает пост со всеми его комментариями.
	// Оптимизированный запрос для получения полной информации о посте.
	GetPostWithComments(ctx context.Context, id uuid.UUID) (*model.PostWithComments, error)

	// Health and lifecycle management

	// HealthCheck проверяет состояние соединения с хранилищем.
	// Возвращает ошибку если хранилище недоступно.
	HealthCheck(ctx context.Context) error

	// Close закрывает соединение с хранилищем и освобождает ресурсы.
	Close() error
}

// PostRepository специализированный интерфейс для работы с постами.
// Предоставляет более специфичные методы для работы с постами.
type PostRepository interface {
	// Create создает новый пост
	Create(ctx context.Context, post *model.Post) (*model.Post, error)

	// GetByID получает пост по ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.Post, error)

	// GetAll получает все посты с пагинацией
	GetAll(ctx context.Context, limit, offset int) ([]*model.Post, error)

	// Update обновляет существующий пост
	Update(ctx context.Context, post *model.Post) (*model.Post, error)

	// Delete удаляет пост
	Delete(ctx context.Context, id uuid.UUID) error

	// ToggleComments включает/выключает комментирование
	ToggleComments(ctx context.Context, id uuid.UUID, enabled bool) error
}

// CommentRepository специализированный интерфейс для работы с комментариями.
// Предоставляет методы для управления комментариями и их иерархией.
type CommentRepository interface {
	// Create создает новый комментарий
	Create(ctx context.Context, comment *model.Comment) (*model.Comment, error)

	// GetByID получает комментарий по ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.Comment, error)

	// GetByPostID получает все комментарии поста
	GetByPostID(ctx context.Context, postID uuid.UUID) ([]model.Comment, error)

	// GetTree получает дерево комментариев поста
	GetTree(ctx context.Context, postID uuid.UUID) ([]model.CommentTree, error)

	// Delete удаляет комментарий
	Delete(ctx context.Context, id uuid.UUID) error
}

// RepositoryManager управляет всеми репозиториями и предоставляет
// единую точку доступа к различным типам репозиториев.
type RepositoryManager interface {
	// Posts возвращает репозиторий для работы с постами
	Posts() PostRepository

	// Comments возвращает репозиторий для работы с комментариями
	Comments() CommentRepository

	// Storage возвращает общий интерфейс хранилища
	Storage() Storage

	// HealthCheck проверяет состояние всех репозиториев
	HealthCheck(ctx context.Context) error

	// Close закрывает все соединения
	Close() error
}
