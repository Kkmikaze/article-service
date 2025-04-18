package repository

import (
	"article-service/internal/domain/article/v1/entity"
	"article-service/pkg/orm"
	articlev1 "article-service/stubs/article/v1"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostRepository interface {
	GetPosts(ctx context.Context, query *orm.QueryBuilder, postStatus *articlev1.PostStatus) ([]*entity.Post, uint64, error)
	GetPostByID(ctx context.Context, id string) (*entity.Post, error)

	// Internal Repos
	InternalGetPosts(ctx context.Context, query *orm.QueryBuilder) ([]*entity.Post, uint64, error)
	InternalGetPostByID(ctx context.Context, id string) (*entity.Post, error)
	InternalCreatePost(ctx context.Context, post *entity.Post) error
	InternalUpdatePost(ctx context.Context, post *entity.Post) (*entity.Post, error)
	InternalDeletePostByID(ctx context.Context, id string) error
}

type postRepository struct {
	Sql *orm.Provider
	Log *logrus.Logger
}

func (r *postRepository) GetPosts(ctx context.Context, query *orm.QueryBuilder, postStatus *articlev1.PostStatus) ([]*entity.Post, uint64, error) {
	var count int64
	var rows []*entity.Post

	statement := `*`
	result := r.Sql.WithContext(ctx).Model(&entity.Post{}).Preload(clause.Associations).Select(statement)
	result = result.Where("status = ?", postStatus)

	if query.Search != "" {
		keyword := "%" + query.Search + "%"
		result = result.Where("title ILIKE ?", keyword)
	}

	if err := result.Count(&count).Error; err != nil {
		r.Log.Errorf("[post repository][func: GetPosts] Failed to count posts: %s", err.Error())

		return nil, 0, status.Error(codes.Internal, "Internal Server Error.")
	}

	if query.Page > 0 {
		result = result.Limit(query.ItemPerPage).Offset((query.Page - 1) * query.ItemPerPage)
	}

	result = result.Order("created_date desc")

	if err := result.Find(&rows).Error; err != nil {
		r.Log.Errorf("[post repository][func: GetPosts] Failed to get posts: %s", err.Error())

		return nil, 0, status.Error(codes.Internal, "Internal Server Error.")
	}

	return rows, uint64(count), nil
}

func (r *postRepository) GetPostByID(ctx context.Context, id string) (*entity.Post, error) {
	var row entity.Post
	if err := r.Sql.WithContext(ctx).Model(&entity.Post{}).Preload(clause.Associations).Where("id = ?", id).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Warnf("[post repository][func: GetPostByID] Post Not Found: %s", err.Error())

			return nil, status.Error(codes.NotFound, "Post Not Found")
		}

		r.Log.Errorf("[post repository][func: GetPostByID] Failed to get post by id: %s", err.Error())

		return nil, status.Error(codes.Internal, "Internal Server Error.")
	}

	return &row, nil
}

// Internal Repos
func (r *postRepository) InternalGetPosts(ctx context.Context, query *orm.QueryBuilder) ([]*entity.Post, uint64, error) {
	var count int64
	var rows []*entity.Post

	statement := `*`
	result := r.Sql.WithContext(ctx).Model(&entity.Post{}).Preload(clause.Associations).Select(statement)

	if query.Search != "" {
		keyword := "%" + query.Search + "%"
		result = result.Where("title ILIKE ?", keyword)
	}

	if err := result.Count(&count).Error; err != nil {
		r.Log.Errorf("[post repository][func: InternalGetPosts] Failed to count posts: %s", err.Error())

		return nil, 0, status.Error(codes.Internal, "Internal Server Error.")
	}

	if query.Page > 0 {
		result = result.Limit(query.ItemPerPage).Offset((query.Page - 1) * query.ItemPerPage)
	}

	result = result.Order("created_date desc")

	if err := result.Find(&rows).Error; err != nil {
		r.Log.Errorf("[post repository][func: InternalGetPosts] Failed to get posts: %s", err.Error())

		return nil, 0, status.Error(codes.Internal, "Internal Server Error.")
	}

	return rows, uint64(count), nil
}

func (r *postRepository) InternalGetPostByID(ctx context.Context, id string) (*entity.Post, error) {
	var row entity.Post
	if err := r.Sql.WithContext(ctx).Model(&entity.Post{}).Preload(clause.Associations).Where("id = ?", id).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Warnf("[post repository][func: InternalGetPostByID] Post Not Found: %s", err.Error())

			return nil, status.Error(codes.NotFound, "Post Not Found")
		}

		r.Log.Errorf("[post repository][func: InternalGetPostByID] Failed to get post by id: %s", err.Error())

		return nil, status.Error(codes.Internal, "Internal Server Error.")
	}

	return &row, nil
}

func (r *postRepository) InternalGetPostByIDs(ctx context.Context, ids []string) ([]*entity.Post, error) {
	var rows []*entity.Post
	if err := r.Sql.WithContext(ctx).Model(&entity.Post{}).Preload(clause.Associations).Where("id IN ?", ids).Find(&rows).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Warnf("[post repository][func: InternalGetPostByIDs] Post Not Found: %s", err.Error())

			return nil, status.Error(codes.NotFound, "Post Not Found")
		}

		r.Log.Errorf("[post repository][func: InternalGetPostByIDs] Failed to get posts by ids: %s", err.Error())

		return nil, status.Error(codes.Internal, "Internal Server Error.")
	}

	return rows, nil
}

func (r *postRepository) InternalCreatePost(ctx context.Context, post *entity.Post) error {
	return r.Sql.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.Sql.WithContext(ctx).Create(post).Error; err != nil {
			r.Log.Errorf("[post repository][func: InternalCreatePost] Failed to create post: %s", err.Error())

			return status.Error(codes.Internal, "Internal Server Error.")
		}

		r.Log.Infof("[post repository][func: InternalCreatePost] Successfully created post with ID: %s", post.ID)
		return nil
	})
}

func (r *postRepository) InternalUpdatePost(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	var updatedPost entity.Post

	err := r.Sql.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existsPost entity.Post
		if err := tx.Model(&entity.Post{}).Where("id != ?", post.ID).First(&existsPost).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				r.Log.Infof("[post repository][func: InternalUpdatePost] Failed to count post: %s\n", err.Error())
				return status.Error(codes.Internal, "Internal Server Error.")
			}
		}

		if err := tx.Model(&entity.Post{}).Where("id = ?", post.ID).Save(post).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				r.Log.Warnf("[post repository][func: InternalUpdatePost] Post Not Found for ID: %s", post.ID)
				return status.Error(codes.NotFound, "Post Not Found")
			}

			r.Log.Errorf("[post repository][func: InternalUpdatePost] Failed to update post with ID %s: %s", post.ID, err.Error())
			return status.Error(codes.Internal, "Internal Server Error.")
		}

		if err := tx.First(&updatedPost, "id = ?", post.ID).Error; err != nil {
			r.Log.Errorf("[post repository][func: InternalUpdatePost] Failed to retrieve updated post with ID %s: %s", post.ID, err.Error())
			return status.Error(codes.Internal, "Internal Server Error.")
		}

		r.Log.Infof("[post repository][func: InternalUpdatePost] Successfully updated post with ID: %s", post.ID)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &updatedPost, nil
}

func (r *postRepository) InternalDeletePostByID(ctx context.Context, id string) error {
	return r.Sql.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var post entity.Post
		if err := tx.First(&post, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				r.Log.Errorf("[post repository][func: InternalDeletePostByID] Post Not Found: %s", err.Error())

				return status.Error(codes.NotFound, "Post Not Found")
			}
		}

		if err := tx.Delete(&post).Error; err != nil {
			r.Log.Errorf("[post repository][func: InternalDeletePostByID] Failed to delete post: %s", err.Error())

			return status.Error(codes.Internal, "Internal Server Error.")
		}

		return nil
	})
}

func NewPostRepository(
	logger *logrus.Logger,
	sql *orm.Provider,
) PostRepository {
	return &postRepository{
		Sql: sql,
		Log: logger,
	}
}
