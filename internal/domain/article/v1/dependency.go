package v1

import (
	"article-service/internal/domain/article/v1/handler"
	"article-service/internal/domain/article/v1/repository"
	"article-service/internal/domain/article/v1/usecase"
	"article-service/pkg/orm"
	articlev1 "article-service/stubs/article/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func RegisterServiceServer(
	rpcServer *grpc.Server,
	logger *logrus.Logger,
	sql *orm.Provider,
) {
	postRepository := repository.NewPostRepository(logger, sql)

	postUsecase := usecase.NewPostUseCase(
		logger,
		postRepository,
	)

	articleManagementHandler := handler.NewArticleHandler(
		logger,
		postUsecase,
	)

	articlev1.RegisterArticleServiceServer(rpcServer, articleManagementHandler)

	grpc_health_v1.RegisterHealthServer(rpcServer, health.NewServer())
}
