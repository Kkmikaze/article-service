package usecase

import (
	"article-service/internal/domain/article/v1/entity"
	"article-service/internal/domain/article/v1/repository"
	"article-service/pkg/orm"
	articlev1 "article-service/stubs/article/v1"
	"context"
	"github.com/sirupsen/logrus"
)

type PostUseCase interface {
	GetPosts(ctx context.Context, in *articlev1.GetPostsRequest) (*articlev1.GetPostsResponse_Data, error)
	GetPostByID(ctx context.Context, in *articlev1.ParamID) (*articlev1.GetPostByIDResponse_Data, error)

	// Internal Usecase
	InternalGetPosts(ctx context.Context, in *articlev1.InternalGetPostsRequest) (*articlev1.InternalGetPostsResponse_Data, error)
	InternalGetPostByID(ctx context.Context, id string) (*articlev1.InternalGetPostByIDResponse_Data, error)
	InternalCreatePost(ctx context.Context, payload *articlev1.InternalCreatePostRequest) error
	InternalUpdatePost(ctx context.Context, payload *articlev1.InternalUpdatePostRequest) error
	InternalDeletePostByID(ctx context.Context, id string) error
}

type articleUseCase struct {
	Log            *logrus.Logger
	PostRepository repository.PostRepository
}

func (u articleUseCase) GetPosts(ctx context.Context, in *articlev1.GetPostsRequest) (*articlev1.GetPostsResponse_Data, error) {
	var items []*articlev1.GetPostsResponse_PostData

	query := &orm.QueryBuilder{
		Search:      in.Search,
		Page:        int(in.Page),
		ItemPerPage: int(in.ItemPerPage),
	}

	rows, total, err := u.PostRepository.GetPosts(ctx, query, &in.Status)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		items = append(items, &articlev1.GetPostsResponse_PostData{
			Id:          row.ID,
			Title:       row.Title,
			Content:     row.Content,
			Category:    row.Category,
			CreatedDate: row.CreatedDate.Format("02 Jan 2006"),
			UpdatedDate: row.UpdatedDate.Format("02 Jan 2006"),
			Status:      row.Status,
		})
	}

	return &articlev1.GetPostsResponse_Data{
		Items: items,
		Total: total,
	}, nil
}

func (u articleUseCase) GetPostByID(ctx context.Context, in *articlev1.ParamID) (*articlev1.GetPostByIDResponse_Data, error) {
	row, err := u.PostRepository.GetPostByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &articlev1.GetPostByIDResponse_Data{
		Id:          row.ID,
		Title:       row.Title,
		Content:     row.Content,
		Category:    row.Category,
		CreatedDate: row.CreatedDate.Format("02 Jan 2006"),
		UpdatedDate: row.UpdatedDate.Format("02 Jan 2006"),
		Status:      row.Status,
	}, nil
}

func (u articleUseCase) InternalGetPosts(ctx context.Context, in *articlev1.InternalGetPostsRequest) (*articlev1.InternalGetPostsResponse_Data, error) {
	var items []*articlev1.InternalGetPostsResponse_PostData

	query := &orm.QueryBuilder{
		Search:      in.Search,
		Page:        int(in.Page),
		ItemPerPage: int(in.ItemPerPage),
	}

	rows, total, err := u.PostRepository.InternalGetPosts(ctx, query)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		items = append(items, &articlev1.InternalGetPostsResponse_PostData{
			Id:          row.ID,
			Title:       row.Title,
			Content:     row.Content,
			Category:    row.Category,
			CreatedDate: row.CreatedDate.Format("02 Jan 2006"),
			UpdatedDate: row.UpdatedDate.Format("02 Jan 2006"),
			Status:      row.Status,
		})
	}

	return &articlev1.InternalGetPostsResponse_Data{
		Items: items,
		Total: total,
	}, nil
}

func (u articleUseCase) InternalGetPostByID(ctx context.Context, id string) (*articlev1.InternalGetPostByIDResponse_Data, error) {
	row, err := u.PostRepository.InternalGetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &articlev1.InternalGetPostByIDResponse_Data{
		Id:          row.ID,
		Title:       row.Title,
		Content:     row.Content,
		Category:    row.Category,
		CreatedDate: row.CreatedDate.Format("02 Jan 2006"),
		UpdatedDate: row.UpdatedDate.Format("02 Jan 2006"),
		Status:      row.Status,
	}, nil
}

func (u articleUseCase) InternalCreatePost(ctx context.Context, payload *articlev1.InternalCreatePostRequest) error {
	article := &entity.Post{
		Title:    payload.Title,
		Content:  payload.Content,
		Category: payload.Category,
		Status:   payload.Status,
	}

	if err := u.PostRepository.InternalCreatePost(ctx, article); err != nil {
		return err
	}

	return nil
}

func (u articleUseCase) InternalUpdatePost(ctx context.Context, payload *articlev1.InternalUpdatePostRequest) error {
	article, err := u.PostRepository.InternalGetPostByID(ctx, payload.Id)
	if err != nil {
		return err
	}

	article.Title = payload.Body.Title
	article.Content = payload.Body.Content
	article.Category = payload.Body.Category
	article.Status = payload.Body.Status

	if _, err := u.PostRepository.InternalUpdatePost(ctx, article); err != nil {
		return err
	}

	return nil
}

func (u articleUseCase) InternalDeletePostByID(ctx context.Context, id string) error {
	if err := u.PostRepository.InternalDeletePostByID(ctx, id); err != nil {
		return err
	}

	return nil
}

func NewPostUseCase(
	logger *logrus.Logger,
	articleRepo repository.PostRepository,
) PostUseCase {
	return &articleUseCase{
		Log:            logger,
		PostRepository: articleRepo,
	}
}
