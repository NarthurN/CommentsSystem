package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/NarthurN/CommentsSystem/internal/model"
	"github.com/NarthurN/CommentsSystem/internal/repository/converter"
	repoModel "github.com/NarthurN/CommentsSystem/internal/repository/model"
)

// PostgresStorage реализует интерфейс Storage для PostgreSQL
type PostgresStorage struct {
	db               *pgxpool.Pool
	postConverter    *converter.PostConverter
	commentConverter *converter.CommentConverter
	treeConverter    *converter.TreeConverter
}

// NewPostgresStorage создает новый экземпляр PostgresStorage
func NewPostgresStorage(ctx context.Context, dsn string) (*PostgresStorage, error) {
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Проверяем соединение
	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresStorage{
		db:               db,
		postConverter:    converter.NewPostConverter(),
		commentConverter: converter.NewCommentConverter(),
		treeConverter:    converter.NewTreeConverter(),
	}, nil
}

// Close закрывает соединение с базой данных
func (s *PostgresStorage) Close() error {
	s.db.Close()
	return nil
}

// HealthCheck проверяет состояние подключения к базе данных
func (s *PostgresStorage) HealthCheck(ctx context.Context) error {
	return s.db.Ping(ctx)
}

// Операции с постами

// CreatePost создает новый пост
func (s *PostgresStorage) CreatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	// Генерируем ID и время создания если не заданы
	if post.ID == uuid.Nil {
		post.ID = uuid.New()
	}
	if post.CreatedAt.IsZero() {
		post.CreatedAt = time.Now()
	}

	// Конвертируем в модель репозитория
	postDB := s.postConverter.ToRepositoryModel(post)

	// Валидируем модель репозитория
	if err := postDB.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Выполняем INSERT
	query := `
		INSERT INTO posts (id, title, content, comments_enabled, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, content, comments_enabled, created_at
	`

	var result repoModel.PostDB
	err := s.db.QueryRow(ctx, query,
		postDB.ID,
		postDB.Title,
		postDB.Content,
		postDB.CommentsEnabled,
		postDB.CreatedAt,
	).Scan(
		&result.ID,
		&result.Title,
		&result.Content,
		&result.CommentsEnabled,
		&result.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Конвертируем обратно в доменную модель
	return s.postConverter.ToDomainModel(&result), nil
}

// GetPost получает пост по ID
func (s *PostgresStorage) GetPost(ctx context.Context, id uuid.UUID) (*model.Post, error) {
	query := `
		SELECT id, title, content, comments_enabled, created_at
		FROM posts
		WHERE id = $1
	`

	var postDB repoModel.PostDB
	err := s.db.QueryRow(ctx, query, id).Scan(
		&postDB.ID,
		&postDB.Title,
		&postDB.Content,
		&postDB.CommentsEnabled,
		&postDB.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return s.postConverter.ToDomainModel(&postDB), nil
}

// GetPosts получает список постов с пагинацией
func (s *PostgresStorage) GetPosts(ctx context.Context, limit, offset int) ([]*model.Post, error) {
	// Значения по умолчанию для пагинации
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT id, title, content, comments_enabled, created_at
		FROM posts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	defer rows.Close()

	var posts []*repoModel.PostDB
	for rows.Next() {
		var postDB repoModel.PostDB
		err := rows.Scan(
			&postDB.ID,
			&postDB.Title,
			&postDB.Content,
			&postDB.CommentsEnabled,
			&postDB.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, &postDB)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	// Конвертируем в доменные модели
	return s.postConverter.ToDomainModels(posts), nil
}

// UpdatePost обновляет пост
func (s *PostgresStorage) UpdatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	postDB := s.postConverter.ToRepositoryModel(post)

	query := `
		UPDATE posts
		SET title = $2, content = $3, comments_enabled = $4
		WHERE id = $1
		RETURNING id, title, content, comments_enabled, created_at
	`

	var result repoModel.PostDB
	err := s.db.QueryRow(ctx, query,
		postDB.ID,
		postDB.Title,
		postDB.Content,
		postDB.CommentsEnabled,
	).Scan(
		&result.ID,
		&result.Title,
		&result.Content,
		&result.CommentsEnabled,
		&result.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return s.postConverter.ToDomainModel(&result), nil
}

// DeletePost удаляет пост
func (s *PostgresStorage) DeletePost(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM posts WHERE id = $1`

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("post not found")
	}

	return nil
}

// TogglePostComments включает/отключает комментарии для поста
func (s *PostgresStorage) TogglePostComments(ctx context.Context, id uuid.UUID, enabled bool) error {
	query := `
		UPDATE posts
		SET comments_enabled = $2
		WHERE id = $1
	`

	result, err := s.db.Exec(ctx, query, id, enabled)
	if err != nil {
		return fmt.Errorf("failed to toggle post comments: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("post not found")
	}

	return nil
}

// Comment operations

// CreateComment создает новый комментарий
func (s *PostgresStorage) CreateComment(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	// Генерируем ID и время создания если не заданы
	if comment.ID == uuid.Nil {
		comment.ID = uuid.New()
	}
	if comment.CreatedAt.IsZero() {
		comment.CreatedAt = time.Now()
	}

	// Конвертируем в модель репозитория
	commentDB := s.commentConverter.ToRepositoryModel(comment)

	// Валидируем модель репозитория
	if err := commentDB.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Выполняем INSERT
	query := `
		INSERT INTO comments (id, post_id, parent_id, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, post_id, parent_id, content, created_at
	`

	var result repoModel.CommentDB
	err := s.db.QueryRow(ctx, query,
		commentDB.ID,
		commentDB.PostID,
		commentDB.ParentID,
		commentDB.Content,
		commentDB.CreatedAt,
	).Scan(
		&result.ID,
		&result.PostID,
		&result.ParentID,
		&result.Content,
		&result.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Конвертируем обратно в доменную модель
	return s.commentConverter.ToDomainModel(&result), nil
}

// GetComment получает комментарий по ID
func (s *PostgresStorage) GetComment(ctx context.Context, id uuid.UUID) (*model.Comment, error) {
	query := `
		SELECT id, post_id, parent_id, content, created_at
		FROM comments
		WHERE id = $1
	`

	var commentDB repoModel.CommentDB
	err := s.db.QueryRow(ctx, query, id).Scan(
		&commentDB.ID,
		&commentDB.PostID,
		&commentDB.ParentID,
		&commentDB.Content,
		&commentDB.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	return s.commentConverter.ToDomainModel(&commentDB), nil
}

// GetCommentsByPostID получает все комментарии для поста
func (s *PostgresStorage) GetCommentsByPostID(ctx context.Context, postID uuid.UUID) ([]model.Comment, error) {
	query := `
		SELECT id, post_id, parent_id, content, created_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC
	`

	rows, err := s.db.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []*repoModel.CommentDB
	for rows.Next() {
		var commentDB repoModel.CommentDB
		err := rows.Scan(
			&commentDB.ID,
			&commentDB.PostID,
			&commentDB.ParentID,
			&commentDB.Content,
			&commentDB.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, &commentDB)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	// Конвертируем в доменные модели
	domainComments := s.commentConverter.ToDomainModels(comments)

	// Конвертируем указатели в значения
	result := make([]model.Comment, len(domainComments))
	for i, comment := range domainComments {
		result[i] = *comment
	}

	return result, nil
}

// GetCommentTree получает иерархическую структуру комментариев для поста
func (s *PostgresStorage) GetCommentTree(ctx context.Context, postID uuid.UUID) ([]model.CommentTree, error) {
	query := `
		WITH RECURSIVE comment_tree AS (
			-- Базовый случай: корневые комментарии
			SELECT id, post_id, parent_id, content, created_at, 0 as level
			FROM comments
			WHERE post_id = $1 AND parent_id IS NULL

			UNION ALL

			-- Рекурсивная часть: дочерние комментарии
			SELECT c.id, c.post_id, c.parent_id, c.content, c.created_at, ct.level + 1
			FROM comments c
			INNER JOIN comment_tree ct ON c.parent_id = ct.id
		)
		SELECT id, post_id, parent_id, content, created_at, level
		FROM comment_tree
		ORDER BY level, created_at
	`

	rows, err := s.db.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment tree: %w", err)
	}
	defer rows.Close()

	var comments []*repoModel.CommentTreeDB
	for rows.Next() {
		var commentDB repoModel.CommentTreeDB
		err := rows.Scan(
			&commentDB.ID,
			&commentDB.PostID,
			&commentDB.ParentID,
			&commentDB.Content,
			&commentDB.CreatedAt,
			&commentDB.Level,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment tree: %w", err)
		}
		comments = append(comments, &commentDB)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	// Строим дерево комментариев
	return s.treeConverter.BuildCommentTree(comments), nil
}

// DeleteComment удаляет комментарий
func (s *PostgresStorage) DeleteComment(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM comments WHERE id = $1`

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}

// Complex operations

// GetPostWithComments получает пост с комментариями
func (s *PostgresStorage) GetPostWithComments(ctx context.Context, id uuid.UUID) (*model.PostWithComments, error) {
	query := `
		SELECT
			p.id, p.title, p.content, p.comments_enabled, p.created_at,
			c.id as comment_id, c.post_id as comment_post_id,
			c.parent_id as comment_parent_id, c.content as comment_content,
			c.created_at as comment_created_at
		FROM posts p
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE p.id = $1
		ORDER BY c.created_at ASC
	`

	rows, err := s.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post with comments: %w", err)
	}
	defer rows.Close()

	var results []*repoModel.PostWithCommentsDB
	for rows.Next() {
		var result repoModel.PostWithCommentsDB
		err := rows.Scan(
			&result.ID,
			&result.Title,
			&result.Content,
			&result.CommentsEnabled,
			&result.CreatedAt,
			&result.CommentID,
			&result.CommentPostID,
			&result.CommentParentID,
			&result.CommentContent,
			&result.CommentCreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post with comments: %w", err)
		}
		results = append(results, &result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("post not found")
	}

	// Конвертируем результат JOIN запроса в PostWithComments
	return s.treeConverter.ToPostWithComments(results), nil
}
