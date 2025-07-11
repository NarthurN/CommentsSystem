package converter

import (
	"testing"
	"time"

	"github.com/NarthurN/CommentsSystem/internal/config"
	"github.com/NarthurN/CommentsSystem/internal/model"
	"github.com/google/uuid"
)

func TestNewGraphQLConverter(t *testing.T) {
	cfg := &config.Config{
		MaxTitleLength:   255,
		MaxContentLength: 10000,
	}

	converter := NewGraphQLConverter(cfg)

	if converter == nil {
		t.Fatal("NewGraphQLConverter returned nil")
	}

	if converter.config != cfg {
		t.Error("Config not properly set in converter")
	}
}

func TestGraphQLConverter_PostToGraphQL(t *testing.T) {
	cfg := &config.Config{
		MaxTitleLength:   255,
		MaxContentLength: 10000,
	}
	converter := NewGraphQLConverter(cfg)

	testTime := time.Now().UTC()
	post := &model.Post{
		ID:              uuid.New(),
		Title:           "Тестовый заголовок",
		Content:         "Тестовое содержимое",
		CommentsEnabled: true,
		CreatedAt:       testTime,
	}

	result, err := converter.PostToGraphQL(post)
	if err != nil {
		t.Fatalf("PostToGraphQL() error = %v", err)
	}

	// Проверяем все поля
	if result["id"].(string) != post.ID.String() {
		t.Errorf("ID mismatch: expected %s, got %s", post.ID.String(), result["id"])
	}

	if result["title"].(string) != post.Title {
		t.Errorf("Title mismatch: expected %s, got %s", post.Title, result["title"])
	}

	if result["content"].(string) != post.Content {
		t.Errorf("Content mismatch: expected %s, got %s", post.Content, result["content"])
	}

	if result["commentsEnabled"].(bool) != post.CommentsEnabled {
		t.Errorf("CommentsEnabled mismatch: expected %v, got %v", post.CommentsEnabled, result["commentsEnabled"])
	}

	expectedTimeStr := testTime.Format("2006-01-02T15:04:05Z07:00")
	if result["createdAt"].(string) != expectedTimeStr {
		t.Errorf("CreatedAt mismatch: expected %s, got %s", expectedTimeStr, result["createdAt"])
	}
}

func TestGraphQLConverter_PostToGraphQL_NilPost(t *testing.T) {
	cfg := &config.Config{}
	converter := NewGraphQLConverter(cfg)

	_, err := converter.PostToGraphQL(nil)
	if err == nil {
		t.Error("Expected error for nil post")
	}
}

func TestGraphQLConverter_CommentToGraphQL(t *testing.T) {
	cfg := &config.Config{}
	converter := NewGraphQLConverter(cfg)

	testTime := time.Now().UTC()
	parentID := uuid.New()
	comment := &model.Comment{
		ID:        uuid.New(),
		PostID:    uuid.New(),
		ParentID:  &parentID,
		Content:   "Тестовый комментарий",
		CreatedAt: testTime,
	}

	result, err := converter.CommentToGraphQL(comment)
	if err != nil {
		t.Fatalf("CommentToGraphQL() error = %v", err)
	}

	// Проверяем все поля
	if result["id"].(string) != comment.ID.String() {
		t.Errorf("ID mismatch: expected %s, got %s", comment.ID.String(), result["id"])
	}

	if result["postId"].(string) != comment.PostID.String() {
		t.Errorf("PostID mismatch: expected %s, got %s", comment.PostID.String(), result["postId"])
	}

	if result["content"].(string) != comment.Content {
		t.Errorf("Content mismatch: expected %s, got %s", comment.Content, result["content"])
	}

	expectedTimeStr := testTime.Format("2006-01-02T15:04:05Z07:00")
	if result["createdAt"].(string) != expectedTimeStr {
		t.Errorf("CreatedAt mismatch: expected %s, got %s", expectedTimeStr, result["createdAt"])
	}

	// Проверяем ParentID
	if result["parentId"].(string) != parentID.String() {
		t.Errorf("ParentID mismatch: expected %s, got %s", parentID.String(), result["parentId"])
	}
}

func TestGraphQLConverter_CommentToGraphQL_NilParentID(t *testing.T) {
	cfg := &config.Config{}
	converter := NewGraphQLConverter(cfg)

	comment := &model.Comment{
		ID:        uuid.New(),
		PostID:    uuid.New(),
		ParentID:  nil, // Корневой комментарий
		Content:   "Корневой комментарий",
		CreatedAt: time.Now().UTC(),
	}

	result, err := converter.CommentToGraphQL(comment)
	if err != nil {
		t.Fatalf("CommentToGraphQL() error = %v", err)
	}

	// ParentID должен быть nil для корневого комментария
	if result["parentId"] != nil {
		t.Errorf("ParentID should be nil for root comment, got %v", result["parentId"])
	}
}

func TestGraphQLConverter_PostsToGraphQL(t *testing.T) {
	cfg := &config.Config{}
	converter := NewGraphQLConverter(cfg)

	posts := []*model.Post{
		{
			ID:              uuid.New(),
			Title:           "Пост 1",
			Content:         "Содержимое 1",
			CommentsEnabled: true,
			CreatedAt:       time.Now().UTC(),
		},
		{
			ID:              uuid.New(),
			Title:           "Пост 2",
			Content:         "Содержимое 2",
			CommentsEnabled: false,
			CreatedAt:       time.Now().UTC(),
		},
	}

	result, err := converter.PostsToGraphQL(posts)
	if err != nil {
		t.Fatalf("PostsToGraphQL() error = %v", err)
	}

	if len(result) != len(posts) {
		t.Errorf("Expected %d posts, got %d", len(posts), len(result))
	}

	// Проверяем первый пост
	if result[0]["title"].(string) != "Пост 1" {
		t.Errorf("First post title mismatch")
	}

	// Проверяем второй пост
	if result[1]["title"].(string) != "Пост 2" {
		t.Errorf("Second post title mismatch")
	}
}

func TestNewValidationConverter(t *testing.T) {
	cfg := &config.Config{
		MaxTitleLength:   255,
		MaxContentLength: 10000,
		MaxCommentLength: 2000,
	}

	converter := NewValidationConverter(cfg)

	if converter == nil {
		t.Fatal("NewValidationConverter returned nil")
	}

	if converter.config != cfg {
		t.Error("Config not properly set in converter")
	}
}

func TestValidationConverter_ValidateAndConvertCreatePost(t *testing.T) {
	cfg := &config.Config{
		MaxTitleLength:   255,
		MaxContentLength: 10000,
	}
	converter := NewValidationConverter(cfg)

	tests := []struct {
		name      string
		title     string
		content   string
		wantError bool
	}{
		{
			name:      "валидный пост",
			title:     "Тестовый заголовок",
			content:   "Тестовое содержимое",
			wantError: false,
		},
		{
			name:      "пустой заголовок",
			title:     "",
			content:   "Тестовое содержимое",
			wantError: true,
		},
		{
			name:      "слишком длинный заголовок",
			title:     string(make([]byte, 256)), // 256 символов
			content:   "Тестовое содержимое",
			wantError: true,
		},
		{
			name:      "пустое содержимое",
			title:     "Тестовый заголовок",
			content:   "",
			wantError: true,
		},
		{
			name:      "слишком длинное содержимое",
			title:     "Тестовый заголовок",
			content:   string(make([]byte, 10001)), // 10001 символ
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := converter.ValidateAndConvertCreatePost(tt.title, tt.content)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, but got none")
				}
				if post != nil {
					t.Error("Expected nil post when error occurs")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if post == nil {
					t.Error("Expected post, but got nil")
				} else {
					if post.Title != tt.title {
						t.Errorf("Title mismatch: expected %s, got %s", tt.title, post.Title)
					}
					if post.Content != tt.content {
						t.Errorf("Content mismatch: expected %s, got %s", tt.content, post.Content)
					}
					if !post.CommentsEnabled {
						t.Error("CommentsEnabled should be true by default")
					}
				}
			}
		})
	}
}

func TestValidationConverter_ValidateAndConvertCreateComment(t *testing.T) {
	cfg := &config.Config{
		MaxCommentLength: 2000,
	}
	converter := NewValidationConverter(cfg)

	validPostID := uuid.New().String()
	validParentID := uuid.New().String()

	tests := []struct {
		name      string
		postID    string
		parentID  string
		content   string
		wantError bool
	}{
		{
			name:      "валидный корневой комментарий",
			postID:    validPostID,
			parentID:  "",
			content:   "Тестовый комментарий",
			wantError: false,
		},
		{
			name:      "валидный дочерний комментарий",
			postID:    validPostID,
			parentID:  validParentID,
			content:   "Ответ на комментарий",
			wantError: false,
		},
		{
			name:      "невалидный postID",
			postID:    "invalid-uuid",
			parentID:  "",
			content:   "Тестовый комментарий",
			wantError: true,
		},
		{
			name:      "невалидный parentID",
			postID:    validPostID,
			parentID:  "invalid-uuid",
			content:   "Тестовый комментарий",
			wantError: true,
		},
		{
			name:      "пустое содержимое",
			postID:    validPostID,
			parentID:  "",
			content:   "",
			wantError: true,
		},
		{
			name:      "слишком длинное содержимое",
			postID:    validPostID,
			parentID:  "",
			content:   string(make([]byte, 2001)), // 2001 символ
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment, err := converter.ValidateAndConvertCreateComment(tt.postID, tt.parentID, tt.content)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, but got none")
				}
				if comment != nil {
					t.Error("Expected nil comment when error occurs")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if comment == nil {
					t.Error("Expected comment, but got nil")
				} else {
					if comment.PostID.String() != tt.postID {
						t.Errorf("PostID mismatch: expected %s, got %s", tt.postID, comment.PostID.String())
					}
					if comment.Content != tt.content {
						t.Errorf("Content mismatch: expected %s, got %s", tt.content, comment.Content)
					}

					// Проверяем ParentID
					if tt.parentID == "" {
						if comment.ParentID != nil {
							t.Error("Expected nil ParentID for root comment")
						}
					} else {
						if comment.ParentID == nil {
							t.Error("Expected non-nil ParentID for child comment")
						} else if comment.ParentID.String() != tt.parentID {
							t.Errorf("ParentID mismatch: expected %s, got %s", tt.parentID, comment.ParentID.String())
						}
					}
				}
			}
		})
	}
}

func TestValidationConverter_ValidatePaginationParams(t *testing.T) {
	cfg := &config.Config{}
	converter := NewValidationConverter(cfg)

	tests := []struct {
		name         string
		limit        *int
		offset       *int
		defaultLimit int
		wantLimit    int
		wantOffset   int
		wantError    bool
	}{
		{
			name:         "nil параметры",
			limit:        nil,
			offset:       nil,
			defaultLimit: 10,
			wantLimit:    10,
			wantOffset:   0,
			wantError:    false,
		},
		{
			name:         "валидные параметры",
			limit:        intPtr(20),
			offset:       intPtr(5),
			defaultLimit: 10,
			wantLimit:    20,
			wantOffset:   5,
			wantError:    false,
		},
		{
			name:         "отрицательный limit",
			limit:        intPtr(-1),
			offset:       intPtr(0),
			defaultLimit: 10,
			wantLimit:    0,
			wantOffset:   0,
			wantError:    true,
		},
		{
			name:         "отрицательный offset",
			limit:        intPtr(10),
			offset:       intPtr(-1),
			defaultLimit: 10,
			wantLimit:    0,
			wantOffset:   0,
			wantError:    true,
		},
		{
			name:         "слишком большой limit",
			limit:        intPtr(1001),
			offset:       intPtr(0),
			defaultLimit: 10,
			wantLimit:    0,
			wantOffset:   0,
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLimit, gotOffset, err := converter.ValidatePaginationParams(tt.limit, tt.offset, tt.defaultLimit)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if gotLimit != tt.wantLimit {
					t.Errorf("Limit mismatch: expected %d, got %d", tt.wantLimit, gotLimit)
				}
				if gotOffset != tt.wantOffset {
					t.Errorf("Offset mismatch: expected %d, got %d", tt.wantOffset, gotOffset)
				}
			}
		})
	}
}

// Вспомогательная функция для создания указателя на int
func intPtr(i int) *int {
	return &i
}
