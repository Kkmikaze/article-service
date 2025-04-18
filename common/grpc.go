package common

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const GRPC_PREFIX_METHOD = "/article.v1.CMSService/"

func GetGRPCRestrictedMethods() map[string][]string {
	return map[string][]string{}
}

func GetGRPCUnrestrictedMethods() []string {
	return []string{
		GRPC_PREFIX_METHOD + "HealthzCheck",
		GRPC_PREFIX_METHOD + "GetPosts",
		GRPC_PREFIX_METHOD + "GetPostByID",
		GRPC_PREFIX_METHOD + "InternalGetPosts",
		GRPC_PREFIX_METHOD + "InternalCreatePost",
		GRPC_PREFIX_METHOD + "InternalGetPostByID",
		GRPC_PREFIX_METHOD + "InternalUpdatePost",
		GRPC_PREFIX_METHOD + "InternalDeletePostByID",
	}
}

func GetMetaData(ctx context.Context) (metadata.MD, error) {
	incomingContext, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "failed to get metadata")
	}

	return incomingContext, nil
}

func GetValueMetaData(metaData metadata.MD, key string) (string, error) {
	if val, ok := metaData[key]; ok {
		return val[0], nil
	}

	return "", status.Error(codes.Internal, "key not found")
}

func TransformGrpcError(err error, code codes.Code, customMessage string) error {
	s, ok := status.FromError(err)
	if !ok {
		message := customMessage
		if message == "" {
			message = err.Error()
		}
		return status.Error(codes.Internal, message)
	}

	message := customMessage
	if message == "" {
		message = s.Message()
	}

	switch code {
	case codes.Unauthenticated:
		return status.Error(codes.Unauthenticated, message)
	case codes.PermissionDenied:
		return status.Error(codes.PermissionDenied, message)
	case codes.NotFound:
		return status.Error(codes.NotFound, message)
	case codes.InvalidArgument:
		return status.Error(codes.InvalidArgument, message)
	case codes.AlreadyExists:
		return status.Error(codes.AlreadyExists, message)
	case codes.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, message)
	case codes.ResourceExhausted:
		return status.Error(codes.ResourceExhausted, message)
	case codes.FailedPrecondition:
		return status.Error(codes.FailedPrecondition, message)
	case codes.Aborted:
		return status.Error(codes.Aborted, message)
	case codes.OutOfRange:
		return status.Error(codes.OutOfRange, message)
	case codes.Unimplemented:
		return status.Error(codes.Unimplemented, message)
	case codes.Internal:
		return status.Error(codes.Internal, message)
	case codes.Unavailable:
		return status.Error(codes.Unavailable, message)
	case codes.DataLoss:
		return status.Error(codes.DataLoss, message)
	case codes.Canceled:
		return status.Error(codes.Canceled, message)
	default:
		// If the code is not explicitly handled, use the original status code from the error
		return status.Error(s.Code(), message)
	}
}
