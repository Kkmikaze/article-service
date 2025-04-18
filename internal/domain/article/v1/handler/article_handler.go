package handler

import (
	"article-service/common"
	"article-service/internal/domain/article/v1/schema"
	"article-service/internal/domain/article/v1/usecase"
	articlev1 "article-service/stubs/article/v1"
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

func (h *articleHandler) HealthzCheck(context.Context, *emptypb.Empty) (*articlev1.HealthCheckResponse, error) {
	return &articlev1.HealthCheckResponse{
		Message: "CMS Service is running.",
	}, nil
}

func (h articleHandler) GetPosts(ctx context.Context, in *articlev1.GetPostsRequest) (*articlev1.GetPostsResponse, error) {
	data, err := h.PostUseCase.GetPosts(ctx, in)
	if err != nil {
		return nil, err
	}

	return &articlev1.GetPostsResponse{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "Get Posts Successfully",
		Data:    data,
	}, nil
}

func (h *articleHandler) GetPostByID(ctx context.Context, in *articlev1.ParamID) (*articlev1.GetPostByIDResponse, error) {
	row, err := h.PostUseCase.GetPostByID(ctx, in)
	if err != nil {
		return nil, err
	}

	return &articlev1.GetPostByIDResponse{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "Get Post by id Successfully.",
		Data:    row,
	}, nil
}

func (h *articleHandler) InternalCreatePost(ctx context.Context, in *articlev1.InternalCreatePostRequest) (*articlev1.CommonResponse, error) {
	validateErr := common.ValidateRequest(&schema.InternalCreatePostRequest{
		Title:    in.Title,
		Content:  in.Content,
		Category: in.Category,
		Status:   in.Status.String(),
	})
	if validateErr != nil {
		return nil, validateErr
	}

	if err := h.PostUseCase.InternalCreatePost(ctx, in); err != nil {
		return nil, err
	}

	return &articlev1.CommonResponse{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "Create Post Successfully.",
	}, nil
}

func (h articleHandler) InternalGetPosts(ctx context.Context, in *articlev1.InternalGetPostsRequest) (*articlev1.InternalGetPostsResponse, error) {
	data, err := h.PostUseCase.InternalGetPosts(ctx, in)
	if err != nil {
		return nil, err
	}

	return &articlev1.InternalGetPostsResponse{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "Get Posts Successfully",
		Data:    data,
	}, nil
}

func (h *articleHandler) InternalGetPostByID(ctx context.Context, in *articlev1.ParamID) (*articlev1.InternalGetPostByIDResponse, error) {
	row, err := h.PostUseCase.InternalGetPostByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &articlev1.InternalGetPostByIDResponse{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "Get Post by id Successfully.",
		Data:    row,
	}, nil
}

func (h *articleHandler) InternalUpdatePost(ctx context.Context, in *articlev1.InternalUpdatePostRequest) (*articlev1.CommonResponse, error) {
	validateErr := common.ValidateRequest(&schema.InternalUpdatePostRequest{
		ID:       in.Id,
		Title:    in.Body.Title,
		Content:  in.Body.Content,
		Category: in.Body.Category,
		Status:   in.Body.Status.String(),
	})
	if validateErr != nil {
		return nil, validateErr
	}

	if err := h.PostUseCase.InternalUpdatePost(ctx, in); err != nil {
		return nil, err
	}

	return &articlev1.CommonResponse{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "Update Post Successfully.",
	}, nil
}

func (h *articleHandler) InternalDeletePostByID(ctx context.Context, in *articlev1.ParamID) (*articlev1.CommonResponse, error) {
	if err := h.PostUseCase.InternalDeletePostByID(ctx, in.Id); err != nil {
		return nil, err
	}

	return &articlev1.CommonResponse{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "Delete Post Successfully.",
	}, nil
}

type articleHandler struct {
	articlev1.UnimplementedArticleServiceServer
	PostUseCase usecase.PostUseCase
	Log         *logrus.Logger
}

func NewArticleHandler(
	logger *logrus.Logger,
	post usecase.PostUseCase,
) articlev1.ArticleServiceServer {
	return &articleHandler{
		PostUseCase: post,
		Log:         logger,
	}
}
