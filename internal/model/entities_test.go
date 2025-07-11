package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPost_IsValidTitle(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected bool
	}{
		{
			name:     "валидный заголовок",
			title:    "Тестовый заголовок",
			expected: true,
		},
		{
			name:     "пустой заголовок",
			title:    "",
			expected: false,
		},
		{
			name:     "очень длинный заголовок",
			title:    string(make([]byte, 256)), // 256 символов
			expected: false,
		},
		{
			name:     "граничный случай - 255 символов",
			title:    string(make([]byte, 255)),
			expected: true,
		},
		{
			name:     "граничный случай - 1 символ",
			title:    "a",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &Post{Title: tt.title}
			result := post.IsValidTitle()
			if result != tt.expected {
				t.Errorf("IsValidTitle() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestPost_IsValidContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "валидное содержимое",
			content:  "Тестовое содержимое поста",
			expected: true,
		},
		{
			name:     "пустое содержимое",
			content:  "",
			expected: false,
		},
		{
			name:     "очень длинное содержимое",
			content:  string(make([]byte, 10001)), // 10001 символ
			expected: false,
		},
		{
			name:     "граничный случай - 10000 символов",
			content:  string(make([]byte, 10000)),
			expected: true,
		},
		{
			name:     "граничный случай - 1 символ",
			content:  "a",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &Post{Content: tt.content}
			result := post.IsValidContent()
			if result != tt.expected {
				t.Errorf("IsValidContent() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestPost_CanAddComments(t *testing.T) {
	tests := []struct {
		name            string
		commentsEnabled bool
		expected        bool
	}{
		{
			name:            "комментарии разрешены",
			commentsEnabled: true,
			expected:        true,
		},
		{
			name:            "комментарии запрещены",
			commentsEnabled: false,
			expected:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &Post{CommentsEnabled: tt.commentsEnabled}
			result := post.CanAddComments()
			if result != tt.expected {
				t.Errorf("CanAddComments() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestPost_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		post     *Post
		expected bool
	}{
		{
			name: "валидный пост",
			post: &Post{
				Title:   "Тестовый заголовок",
				Content: "Тестовое содержимое",
			},
			expected: true,
		},
		{
			name: "невалидный заголовок",
			post: &Post{
				Title:   "",
				Content: "Тестовое содержимое",
			},
			expected: false,
		},
		{
			name: "невалидное содержимое",
			post: &Post{
				Title:   "Тестовый заголовок",
				Content: "",
			},
			expected: false,
		},
		{
			name: "оба поля невалидны",
			post: &Post{
				Title:   "",
				Content: "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.post.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestPost_Prepare(t *testing.T) {
	t.Run("подготовка нового поста", func(t *testing.T) {
		post := &Post{
			Title:   "Тестовый заголовок",
			Content: "Тестовое содержимое",
		}

		// До вызова Prepare()
		if post.ID != uuid.Nil {
			t.Error("ID должен быть пустым до вызова Prepare()")
		}
		if !post.CreatedAt.IsZero() {
			t.Error("CreatedAt должен быть пустым до вызова Prepare()")
		}

		post.Prepare()

		// После вызова Prepare()
		if post.ID == uuid.Nil {
			t.Error("ID должен быть сгенерирован после вызова Prepare()")
		}
		if post.CreatedAt.IsZero() {
			t.Error("CreatedAt должен быть установлен после вызова Prepare()")
		}
		if time.Since(post.CreatedAt) > time.Second {
			t.Error("CreatedAt должен быть близок к текущему времени")
		}
	})

	t.Run("не перезаписывает существующие ID и время", func(t *testing.T) {
		existingID := uuid.New()
		existingTime := time.Now().Add(-time.Hour)

		post := &Post{
			ID:        existingID,
			CreatedAt: existingTime,
			Title:     "Тестовый заголовок",
			Content:   "Тестовое содержимое",
		}

		post.Prepare()

		if post.ID != existingID {
			t.Error("Существующий ID не должен быть перезаписан")
		}
		if post.CreatedAt != existingTime {
			t.Error("Существующее время не должно быть перезаписано")
		}
	})
}

func TestComment_IsValidComment(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "валидный комментарий",
			content:  "Тестовый комментарий",
			expected: true,
		},
		{
			name:     "пустой комментарий",
			content:  "",
			expected: false,
		},
		{
			name:     "очень длинный комментарий",
			content:  string(make([]byte, 2001)), // 2001 символ
			expected: false,
		},
		{
			name:     "граничный случай - 2000 символов",
			content:  string(make([]byte, 2000)),
			expected: true,
		},
		{
			name:     "граничный случай - 1 символ",
			content:  "a",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment := &Comment{Content: tt.content}
			result := comment.IsValidComment()
			if result != tt.expected {
				t.Errorf("IsValidComment() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestComment_IsRootComment(t *testing.T) {
	tests := []struct {
		name     string
		parentID *uuid.UUID
		expected bool
	}{
		{
			name:     "корневой комментарий",
			parentID: nil,
			expected: true,
		},
		{
			name:     "дочерний комментарий",
			parentID: func() *uuid.UUID { id := uuid.New(); return &id }(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment := &Comment{ParentID: tt.parentID}
			result := comment.IsRootComment()
			if result != tt.expected {
				t.Errorf("IsRootComment() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestComment_HasValidPost(t *testing.T) {
	tests := []struct {
		name     string
		postID   uuid.UUID
		expected bool
	}{
		{
			name:     "валидный PostID",
			postID:   uuid.New(),
			expected: true,
		},
		{
			name:     "пустой PostID",
			postID:   uuid.Nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment := &Comment{PostID: tt.postID}
			result := comment.HasValidPost()
			if result != tt.expected {
				t.Errorf("HasValidPost() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestComment_IsValid(t *testing.T) {
	validPostID := uuid.New()

	tests := []struct {
		name     string
		comment  *Comment
		expected bool
	}{
		{
			name: "валидный комментарий",
			comment: &Comment{
				PostID:  validPostID,
				Content: "Тестовый комментарий",
			},
			expected: true,
		},
		{
			name: "невалидное содержимое",
			comment: &Comment{
				PostID:  validPostID,
				Content: "",
			},
			expected: false,
		},
		{
			name: "невалидный PostID",
			comment: &Comment{
				PostID:  uuid.Nil,
				Content: "Тестовый комментарий",
			},
			expected: false,
		},
		{
			name: "оба поля невалидны",
			comment: &Comment{
				PostID:  uuid.Nil,
				Content: "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.comment.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestComment_Prepare(t *testing.T) {
	t.Run("подготовка нового комментария", func(t *testing.T) {
		comment := &Comment{
			PostID:  uuid.New(),
			Content: "Тестовый комментарий",
		}

		// До вызова Prepare()
		if comment.ID != uuid.Nil {
			t.Error("ID должен быть пустым до вызова Prepare()")
		}
		if !comment.CreatedAt.IsZero() {
			t.Error("CreatedAt должен быть пустым до вызова Prepare()")
		}

		comment.Prepare()

		// После вызова Prepare()
		if comment.ID == uuid.Nil {
			t.Error("ID должен быть сгенерирован после вызова Prepare()")
		}
		if comment.CreatedAt.IsZero() {
			t.Error("CreatedAt должен быть установлен после вызова Prepare()")
		}
		if time.Since(comment.CreatedAt) > time.Second {
			t.Error("CreatedAt должен быть близок к текущему времени")
		}
	})
}

func TestCommentTree_HasChildren(t *testing.T) {
	tests := []struct {
		name     string
		children []CommentTree
		expected bool
	}{
		{
			name:     "есть дочерние комментарии",
			children: []CommentTree{{}},
			expected: true,
		},
		{
			name:     "нет дочерних комментариев",
			children: []CommentTree{},
			expected: false,
		},
		{
			name:     "nil дочерние комментарии",
			children: nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commentTree := &CommentTree{Children: tt.children}
			result := commentTree.HasChildren()
			if result != tt.expected {
				t.Errorf("HasChildren() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCommentTree_GetChildrenCount(t *testing.T) {
	tests := []struct {
		name          string
		children      []CommentTree
		expectedCount int
	}{
		{
			name:          "нет дочерних комментариев",
			children:      []CommentTree{},
			expectedCount: 0,
		},
		{
			name:          "один дочерний комментарий",
			children:      []CommentTree{{}},
			expectedCount: 1,
		},
		{
			name:          "несколько дочерних комментариев",
			children:      []CommentTree{{}, {}, {}},
			expectedCount: 3,
		},
		{
			name:          "nil дочерние комментарии",
			children:      nil,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commentTree := &CommentTree{Children: tt.children}
			result := commentTree.GetChildrenCount()
			if result != tt.expectedCount {
				t.Errorf("GetChildrenCount() = %v, expected %v", result, tt.expectedCount)
			}
		})
	}
}

// Benchmark тесты для производительности
func BenchmarkPost_IsValid(b *testing.B) {
	post := &Post{
		Title:   "Тестовый заголовок для бенчмарка",
		Content: "Тестовое содержимое для бенчмарка с достаточным количеством текста",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		post.IsValid()
	}
}

func BenchmarkComment_IsValid(b *testing.B) {
	comment := &Comment{
		PostID:  uuid.New(),
		Content: "Тестовый комментарий для бенчмарка",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		comment.IsValid()
	}
}
