package converter

import (
	"time"

	"github.com/google/uuid"

	"github.com/NarthurN/CommentsSystem/internal/model"
	repoModel "github.com/NarthurN/CommentsSystem/internal/repository/model"
)

// PostConverter отвечает за конвертацию между доменными моделями и моделями репозитория
type PostConverter struct{}

// NewPostConverter создает новый экземпляр PostConverter
func NewPostConverter() *PostConverter {
	return &PostConverter{}
}

// ToRepositoryModel конвертирует доменную модель поста в модель репозитория
func (c *PostConverter) ToRepositoryModel(domainPost *model.Post) *repoModel.PostDB {
	if domainPost == nil {
		return nil
	}

	return &repoModel.PostDB{
		ID:              domainPost.ID,
		Title:           domainPost.Title,
		Content:         domainPost.Content,
		CommentsEnabled: domainPost.CommentsEnabled,
		CreatedAt:       domainPost.CreatedAt,
	}
}

// ToDomainModel конвертирует модель репозитория поста в доменную модель
func (c *PostConverter) ToDomainModel(repoPost *repoModel.PostDB) *model.Post {
	if repoPost == nil {
		return nil
	}

	return &model.Post{
		ID:              repoPost.ID,
		Title:           repoPost.Title,
		Content:         repoPost.Content,
		CommentsEnabled: repoPost.CommentsEnabled,
		CreatedAt:       repoPost.CreatedAt,
	}
}

// ToDomainModels конвертирует слайс моделей репозитория в слайс доменных моделей
func (c *PostConverter) ToDomainModels(repoPosts []*repoModel.PostDB) []*model.Post {
	if repoPosts == nil {
		return nil
	}

	domainPosts := make([]*model.Post, len(repoPosts))
	for i, repoPost := range repoPosts {
		domainPosts[i] = c.ToDomainModel(repoPost)
	}

	return domainPosts
}

// CreateNewPost создает новую доменную модель поста с сгенерированным ID
func (c *PostConverter) CreateNewPost(title, content string, commentsEnabled bool) *model.Post {
	return &model.Post{
		ID:              uuid.New(),
		Title:           title,
		Content:         content,
		CommentsEnabled: commentsEnabled,
		CreatedAt:       time.Now(),
	}
}

// CommentConverter отвечает за конвертацию между доменными моделями и моделями репозитория
type CommentConverter struct{}

// NewCommentConverter создает новый экземпляр CommentConverter
func NewCommentConverter() *CommentConverter {
	return &CommentConverter{}
}

// ToRepositoryModel конвертирует доменную модель комментария в модель репозитория
func (c *CommentConverter) ToRepositoryModel(domainComment *model.Comment) *repoModel.CommentDB {
	if domainComment == nil {
		return nil
	}

	return &repoModel.CommentDB{
		ID:        domainComment.ID,
		PostID:    domainComment.PostID,
		ParentID:  domainComment.ParentID,
		Content:   domainComment.Content,
		CreatedAt: domainComment.CreatedAt,
	}
}

// ToDomainModel конвертирует модель репозитория комментария в доменную модель
func (c *CommentConverter) ToDomainModel(repoComment *repoModel.CommentDB) *model.Comment {
	if repoComment == nil {
		return nil
	}

	return &model.Comment{
		ID:        repoComment.ID,
		PostID:    repoComment.PostID,
		ParentID:  repoComment.ParentID,
		Content:   repoComment.Content,
		CreatedAt: repoComment.CreatedAt,
	}
}

// ToDomainModels конвертирует слайс моделей репозитория в слайс доменных моделей
func (c *CommentConverter) ToDomainModels(repoComments []*repoModel.CommentDB) []*model.Comment {
	if repoComments == nil {
		return nil
	}

	domainComments := make([]*model.Comment, len(repoComments))
	for i, repoComment := range repoComments {
		domainComments[i] = c.ToDomainModel(repoComment)
	}

	return domainComments
}

// CreateNewComment создает новую доменную модель комментария с сгенерированным ID
func (c *CommentConverter) CreateNewComment(postID uuid.UUID, parentID *uuid.UUID, content string) *model.Comment {
	return &model.Comment{
		ID:        uuid.New(),
		PostID:    postID,
		ParentID:  parentID,
		Content:   content,
		CreatedAt: time.Now(),
	}
}

// TreeConverter отвечает за конвертацию древовидных структур комментариев
type TreeConverter struct {
	commentConverter *CommentConverter
}

// NewTreeConverter создает новый экземпляр TreeConverter
func NewTreeConverter() *TreeConverter {
	return &TreeConverter{
		commentConverter: NewCommentConverter(),
	}
}

// BuildCommentTree строит иерархическую структуру комментариев из плоского списка
func (c *TreeConverter) BuildCommentTree(repoComments []*repoModel.CommentTreeDB) []model.CommentTree {
	if len(repoComments) == 0 {
		return nil
	}

	// Конвертируем в доменные модели
	domainComments := make([]model.Comment, len(repoComments))
	for i, repoComment := range repoComments {
		domainComments[i] = *c.commentConverter.ToDomainModel(&repoComment.CommentDB)
	}

	// Строим дерево
	return c.buildTree(domainComments, nil)
}

// buildTree рекурсивно строит дерево комментариев
func (c *TreeConverter) buildTree(comments []model.Comment, parentID *uuid.UUID) []model.CommentTree {
	var result []model.CommentTree

	for _, comment := range comments {
		// Проверяем, является ли этот комментарий дочерним для текущего parentID
		if (parentID == nil && comment.ParentID == nil) ||
			(parentID != nil && comment.ParentID != nil && *comment.ParentID == *parentID) {

			tree := model.CommentTree{
				Comment:  comment,
				Children: c.buildTree(comments, &comment.ID),
			}
			result = append(result, tree)
		}
	}

	return result
}

// ToPostWithComments конвертирует результат JOIN запроса в PostWithComments
func (c *TreeConverter) ToPostWithComments(repoResult []*repoModel.PostWithCommentsDB) *model.PostWithComments {
	if len(repoResult) == 0 {
		return nil
	}

	// Берем информацию о посте из первой записи
	firstRow := repoResult[0]
	postConverter := NewPostConverter()
	post := postConverter.ToDomainModel(&firstRow.PostDB)

	// Собираем уникальные комментарии
	commentMap := make(map[uuid.UUID]*model.Comment)
	for _, row := range repoResult {
		if row.CommentID != nil {
			comment := &model.Comment{
				ID:        *row.CommentID,
				PostID:    *row.CommentPostID,
				ParentID:  row.CommentParentID,
				Content:   *row.CommentContent,
				CreatedAt: *row.CommentCreatedAt,
			}
			commentMap[comment.ID] = comment
		}
	}

	// Конвертируем в слайс
	comments := make([]model.Comment, 0, len(commentMap))
	for _, comment := range commentMap {
		comments = append(comments, *comment)
	}

	return &model.PostWithComments{
		Post:     *post,
		Comments: comments,
	}
}
