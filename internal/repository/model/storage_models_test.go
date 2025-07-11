package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPostDB_Validate(t *testing.T) {
	tests := []struct {
		name      string
		post      PostDB
		wantError bool
	}{
		{
			name: "валидный пост",
			post: PostDB{
				ID:              uuid.New(),
				Title:           "Тестовый заголовок",
				Content:         "Тестовое содержимое",
				CommentsEnabled: true,
				CreatedAt:       time.Now().UTC(),
			},
			wantError: false,
		},
		{
			name: "пустой заголовок",
			post: PostDB{
				ID:              uuid.New(),
				Title:           "",
				Content:         "Тестовое содержимое",
				CommentsEnabled: true,
				CreatedAt:       time.Now().UTC(),
			},
			wantError: true,
		},
		{
			name: "пустое содержимое",
			post: PostDB{
				ID:              uuid.New(),
				Title:           "Тестовый заголовок",
				Content:         "",
				CommentsEnabled: true,
				CreatedAt:       time.Now().UTC(),
			},
			wantError: true,
		},
		{
			name: "пустой ID",
			post: PostDB{
				ID:              uuid.Nil,
				Title:           "Тестовый заголовок",
				Content:         "Тестовое содержимое",
				CommentsEnabled: true,
				CreatedAt:       time.Now().UTC(),
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.post.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("PostDB.Validate() error = %v, wantErr %v", err, tt.wantError)
			}
		})
	}
}

func TestCommentDB_Validate(t *testing.T) {
	tests := []struct {
		name      string
		comment   CommentDB
		wantError bool
	}{
		{
			name: "валидный корневой комментарий",
			comment: CommentDB{
				ID:        uuid.New(),
				PostID:    uuid.New(),
				ParentID:  nil,
				Content:   "Тестовый комментарий",
				CreatedAt: time.Now().UTC(),
			},
			wantError: false,
		},
		{
			name: "валидный дочерний комментарий",
			comment: CommentDB{
				ID:        uuid.New(),
				PostID:    uuid.New(),
				ParentID:  func() *uuid.UUID { id := uuid.New(); return &id }(),
				Content:   "Тестовый комментарий",
				CreatedAt: time.Now().UTC(),
			},
			wantError: false,
		},
		{
			name: "пустое содержимое",
			comment: CommentDB{
				ID:        uuid.New(),
				PostID:    uuid.New(),
				ParentID:  nil,
				Content:   "",
				CreatedAt: time.Now().UTC(),
			},
			wantError: true,
		},
		{
			name: "пустой PostID",
			comment: CommentDB{
				ID:        uuid.New(),
				PostID:    uuid.Nil,
				ParentID:  nil,
				Content:   "Тестовый комментарий",
				CreatedAt: time.Now().UTC(),
			},
			wantError: true,
		},
		{
			name: "пустой ID",
			comment: CommentDB{
				ID:        uuid.Nil,
				PostID:    uuid.New(),
				ParentID:  nil,
				Content:   "Тестовый комментарий",
				CreatedAt: time.Now().UTC(),
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.comment.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("CommentDB.Validate() error = %v, wantErr %v", err, tt.wantError)
			}
		})
	}
}
