package interceptors

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ServerInterceptor struct {
	logger             *logrus.Logger
	restrictedMethod   map[string][]string
	unrestrictedMethod map[string]struct{}
	publicKey          []byte
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (ss *wrappedStream) Context() context.Context {
	return ss.ctx
}

func (ss *wrappedStream) SetContext(ctx context.Context) {
	ss.ctx = ctx
}

func (ss *wrappedStream) RecvMsg(m interface{}) error {
	return ss.ServerStream.RecvMsg(m)
}

func (ss *wrappedStream) SendMsg(m interface{}) error {
	return ss.ServerStream.SendMsg(m)
}

type StreamContextWrapper interface {
	grpc.ServerStream
	SetContext(context.Context)
}

func newStreamContextWrapper(ss grpc.ServerStream, ctx context.Context) StreamContextWrapper {
	return &wrappedStream{
		ss,
		ctx,
	}
}

type ServerInterceptorInterface interface {
	IsRestrictedMethodAllowed(method string, accesses []string) bool
	RegisterRestrictedMethods(methods map[string][]string)
	RegisterUnrestrictedMethods(methods []string)
	UnaryServerMetadataPropagationInterceptor() grpc.UnaryServerInterceptor
	UnaryServerRecoveryInterceptor() grpc.UnaryServerInterceptor
	UnaryServerAuthInterceptor(publicKey []byte) grpc.UnaryServerInterceptor
	UnaryServerPerformanceInterceptor(customLog ServerInterceptorPerformanceLog) grpc.UnaryServerInterceptor
	ServerStreamMetadataPropagationInterceptor() grpc.StreamServerInterceptor
	ServerStreamRecoveryInterceptor() grpc.StreamServerInterceptor
	ServerStreamAuthInterceptor(publicKey []byte) grpc.StreamServerInterceptor
	ServerStreamPerformanceInterceptor(customLog ServerInterceptorPerformanceLog) grpc.StreamServerInterceptor
}

type ServerInterceptorPerformanceLog interface {
	Infof(format string, args ...interface{})
}

func (s *ServerInterceptor) RegisterRestrictedMethods(methods map[string][]string) {
	s.restrictedMethod = methods
}

// RegisterUnrestrictedMethods allows setting unrestricted methods.
func (s *ServerInterceptor) RegisterUnrestrictedMethods(methods []string) {
	s.unrestrictedMethod = make(map[string]struct{}, len(methods))
	for _, method := range methods {
		s.unrestrictedMethod[method] = struct{}{}
	}
}

func (s *ServerInterceptor) IsRestrictedMethodAllowed(method string, accesses []string) bool {
	rules, ok := s.restrictedMethod[method]
	if !ok {
		return true
	}

	ruleMap := make(map[string]struct{}, len(rules))
	for _, rule := range rules {
		ruleMap[rule] = struct{}{}
	}

	for _, access := range accesses {
		if _, exists := ruleMap[access]; exists {
			return true
		}
	}
	return false
}

func (s *ServerInterceptor) UnaryServerRecoveryInterceptor() grpc.UnaryServerInterceptor {
	return recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p any) (err error) {
		// Log the error message with type and additional context
		logrus.Errorf("Panic recovered in gRPC stream: %v (type: %T)", p, p)

		// Optionally log the stack trace for debugging (if stack trace collection is supported)
		logrus.Errorf("Stack trace: %s", debug.Stack())

		// Return a structured error to the client
		return status.Errorf(codes.Unknown, "Unknown Server Error")
	}))
}

func (s *ServerInterceptor) UnaryServerAuthInterceptor(publicKey []byte) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if _, ok := s.unrestrictedMethod[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		return handler(ctx, req)
	}
}

func (s *ServerInterceptor) UnaryServerPerformanceInterceptor(customLog ServerInterceptorPerformanceLog) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)

		statusCode := codes.OK
		if err != nil {
			if statusErr, ok := status.FromError(err); ok {
				statusCode = statusErr.Code()
			} else {
				statusCode = codes.Internal
			}
		}

		duration := time.Since(start)
		customLog.Infof("[Latency: %s, Method: %s, StatusCode: %v]", duration, info.FullMethod, statusCode)

		return resp, err
	}
}

func (s *ServerInterceptor) UnaryServerMetadataPropagationInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		incomingMD, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			incomingMD = metadata.MD{}
		}

		outgoingCtx := metadata.NewOutgoingContext(ctx, incomingMD)

		return handler(outgoingCtx, req)
	}
}

func (s *ServerInterceptor) ServerStreamMetadataPropagationInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		incomingMD, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			incomingMD = metadata.MD{}
		}

		outgoingCtx := metadata.NewOutgoingContext(ss.Context(), incomingMD)
		wrapped := newStreamContextWrapper(ss, outgoingCtx)

		return handler(srv, wrapped)
	}
}

func (s *ServerInterceptor) ServerStreamAuthInterceptor(publicKey []byte) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if _, ok := s.unrestrictedMethod[info.FullMethod]; ok {
			return handler(srv, ss)
		}

		wrapped := newStreamContextWrapper(ss, ss.Context())

		return handler(srv, wrapped)
	}
}

func (s *ServerInterceptor) ServerStreamPerformanceInterceptor(customLog ServerInterceptorPerformanceLog) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		err := handler(srv, ss)

		statusCode := codes.OK
		if err != nil {
			if statusErr, ok := status.FromError(err); ok {
				statusCode = statusErr.Code()
			} else {
				statusCode = codes.Internal
			}
		}

		duration := time.Since(start)
		customLog.Infof("[Latency: %s, Method: %s, StatusCode: %v]", duration, info.FullMethod, statusCode)

		return err
	}
}

func (s *ServerInterceptor) ServerStreamRecoveryInterceptor() grpc.StreamServerInterceptor {
	return recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(func(p any) (err error) {
		// Log the error message with type and additional context
		logrus.Errorf("Panic recovered in gRPC stream: %v (type: %T)", p, p)

		// Optionally log the stack trace for debugging (if stack trace collection is supported)
		logrus.Errorf("Stack trace: %s", debug.Stack())

		// Return a structured error to the client
		return status.Errorf(codes.Unknown, "Unknown Server Error")
	}))
}

func NewServerInterceptor(logger *logrus.Logger) ServerInterceptorInterface {
	return &ServerInterceptor{
		logger: logger,
	}
}
