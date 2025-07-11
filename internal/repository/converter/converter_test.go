package converter

import (
	"testing"
	"time"

	"github.com/NarthurN/CommentsSystem/internal/model"
	repoModel "github.com/NarthurN/CommentsSystem/internal/repository/model"
	"github.com/google/uuid"
)

func TestNewPostConverter(t *testing.T) {
	converter := NewPostConverter()
	if converter == nil {
		t.Fatal("NewPostConverter returned nil")
	}
}

func TestPostConverter_ToRepositoryModel(t *testing.T) {
	converter := NewPostConverter()

	testTime := time.Now().UTC()
	domainPost := &model.Post{
		ID:              uuid.New(),
		Title:           "Тестовый заголовок",
		Content:         "Тестовое содержимое",
		CommentsEnabled: true,
		CreatedAt:       testTime,
	}

	repoPost := converter.ToRepositoryModel(domainPost)

	if repoPost.ID != domainPost.ID {
		t.Errorf("ID mismatch: expected %v, got %v", domainPost.ID, repoPost.ID)
	}

	if repoPost.Title != domainPost.Title {
		t.Errorf("Title mismatch: expected %s, got %s", domainPost.Title, repoPost.Title)
	}

	if repoPost.Content != domainPost.Content {
		t.Errorf("Content mismatch: expected %s, got %s", domainPost.Content, repoPost.Content)
	}

	if repoPost.CommentsEnabled != domainPost.CommentsEnabled {
		t.Errorf("CommentsEnabled mismatch: expected %v, got %v", domainPost.CommentsEnabled, repoPost.CommentsEnabled)
	}

	if !repoPost.CreatedAt.Equal(domainPost.CreatedAt) {
		t.Errorf("CreatedAt mismatch: expected %v, got %v", domainPost.CreatedAt, repoPost.CreatedAt)
	}
}

func TestPostConverter_ToDomainModel(t *testing.T) {
	converter := NewPostConverter()

	testTime := time.Now().UTC()
	repoPost := &repoModel.PostDB{
		ID:              uuid.New(),
		Title:           "Тестовый заголовок",
		Content:         "Тестовое содержимое",
		CommentsEnabled: true,
		CreatedAt:       testTime,
	}

	domainPost := converter.ToDomainModel(repoPost)

	if domainPost.ID != repoPost.ID {
		t.Errorf("ID mismatch: expected %v, got %v", repoPost.ID, domainPost.ID)
	}

	if domainPost.Title != repoPost.Title {
		t.Errorf("Title mismatch: expected %s, got %s", repoPost.Title, domainPost.Title)
	}
}

func TestPostConverter_ToDomainModels(t *testing.T) {
	converter := NewPostConverter()

	repoPosts := []*repoModel.PostDB{
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

	domainPosts := converter.ToDomainModels(repoPosts)

	if len(domainPosts) != len(repoPosts) {
		t.Errorf("Length mismatch: expected %d, got %d", len(repoPosts), len(domainPosts))
	}

	for i, domainPost := range domainPosts {
		if domainPost.Title != repoPosts[i].Title {
			t.Errorf("Post %d title mismatch: expected %s, got %s", i, repoPosts[i].Title, domainPost.Title)
		}
	}
}

func TestNewCommentConverter(t *testing.T) {
	converter := NewCommentConverter()
	if converter == nil {
		t.Fatal("NewCommentConverter returned nil")
	}
}

func TestCommentConverter_ToRepositoryModel(t *testing.T) {
	converter := NewCommentConverter()

	testTime := time.Now().UTC()
	parentID := uuid.New()
	domainComment := &model.Comment{
		ID:        uuid.New(),
		PostID:    uuid.New(),
		ParentID:  &parentID,
		Content:   "Тестовый комментарий",
		CreatedAt: testTime,
	}

	repoComment := converter.ToRepositoryModel(domainComment)

	if repoComment.ID != domainComment.ID {
		t.Errorf("ID mismatch: expected %v, got %v", domainComment.ID, repoComment.ID)
	}

	if repoComment.Content != domainComment.Content {
		t.Errorf("Content mismatch: expected %s, got %s", domainComment.Content, repoComment.Content)
	}
}

func TestCommentConverter_ToDomainModel(t *testing.T) {
	converter := NewCommentConverter()

	testTime := time.Now().UTC()
	parentID := uuid.New()
	repoComment := &repoModel.CommentDB{
		ID:        uuid.New(),
		PostID:    uuid.New(),
		ParentID:  &parentID,
		Content:   "Тестовый комментарий",
		CreatedAt: testTime,
	}

	domainComment := converter.ToDomainModel(repoComment)

	if domainComment.ID != repoComment.ID {
		t.Errorf("ID mismatch: expected %v, got %v", repoComment.ID, domainComment.ID)
	}

	if domainComment.Content != repoComment.Content {
		t.Errorf("Content mismatch: expected %s, got %s", repoComment.Content, domainComment.Content)
	}
}

func TestNewTreeConverter(t *testing.T) {
	converter := NewTreeConverter()
	if converter == nil {
		t.Fatal("NewTreeConverter returned nil")
	}
}
