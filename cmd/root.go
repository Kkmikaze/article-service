package cmd

import (
	"article-service/common"
	"article-service/config"
	"article-service/constants"
	articleDependency "article-service/internal/domain/article/v1"
	"article-service/pkg/gateway"
	"article-service/pkg/interceptors"
	"article-service/pkg/orm"
	"article-service/pkg/rpcclient"
	"article-service/pkg/rpcserver"
	articlev1 "article-service/stubs/article/v1"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	rpcPort     string
	gatewayPort string
	rootCmd     = &cobra.Command{
		Use:   "service",
		Short: "Running the gRPC service",
		Long:  "Used to run gRPC Service including rpc server, rpc client, and gateway",
		Run:   runService,
	}
)

func runService(cmd *cobra.Command, args []string) {
	logger := logrus.New()
	ctx := context.Background()

	// Setup interceptors and SSL
	interceptor := interceptors.NewServerInterceptor(logger)
	interceptor.RegisterRestrictedMethods(common.GetGRPCRestrictedMethods())
	interceptor.RegisterUnrestrictedMethods(common.GetGRPCUnrestrictedMethods())

	// Initialize RPC Server
	rpcServer := rpcserver.NewRPCServer(rpcPort,
		"tcp",
		false,
		[]grpc.UnaryServerInterceptor{
			interceptor.UnaryServerRecoveryInterceptor(),
			interceptor.UnaryServerMetadataPropagationInterceptor(),
			interceptor.UnaryServerPerformanceInterceptor(logger),
		},
		[]grpc.StreamServerInterceptor{
			interceptor.ServerStreamRecoveryInterceptor(),
			interceptor.ServerStreamMetadataPropagationInterceptor(),
			interceptor.ServerStreamPerformanceInterceptor(logger),
		},
	)

	// postgresql connection
	sql, err := orm.NewPostgreSQL(ctx,
		config.Configs.PostgresDNS,
		&orm.ConfigConnProvider{
			ConnMaxLifetime: constants.OrmConnMaxLifeTime,
			ConnMaxIdleTime: constants.OrmConnMaxIdleTime,
			MaxOpenConns:    constants.OrmMaxOpenConns,
			MaxIdleConns:    constants.OrmMaxIdleConns,
		},
		&gorm.Config{},
	)
	if err != nil {
		logger.Fatalln(fmt.Errorf("failed to connect to postgresql: %w", err))
	}

	articleDependency.RegisterServiceServer(rpcServer.Server, logger, sql)

	defer rpcServer.StopListener()
	logger.Infoln("Serving gRPC on", rpcPort)

	if err := rpcServer.Run(); err != nil {
		logger.Fatalln("Failed to listen grpc server", err)
	}

	// RPC Client
	client, err := rpcclient.NewRPCClient(rpcPort, []grpc.UnaryClientInterceptor{})
	if err != nil {
		logger.Fatalln("Failed to dial gRPC server", err)
	}

	// API Gateway
	apiGateway := setupAPIGateway()

	err = articlev1.RegisterArticleServiceHandler(ctx, apiGateway.GetRuntimeMux(), client.UnaryConn)
	if err != nil {
		logger.Fatalln("Failed to register user Management service", err)
	}

	logger.Infoln("Serving gRPC-Gateway on", gatewayPort)

	// Enable Swagger
	common.EnableSwagger(apiGateway)

	// Graceful Shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	errs := make(chan error, 2)

	// Run API Gateway
	go func() {
		if err := apiGateway.Run(ctx, common.HandlerMux(apiGateway, []string{"application/json", "application/x-www-form-urlencoded"})); err != nil {
			logger.Fatalln("Failed to listen grpc gateway", err)
			errs <- err
		}
	}()

	// Wait for error or termination signal
	select {
	case err := <-errs:
		logger.Fatalln("Error occurred:", err)
	case sig := <-sigs:
		logger.Infoln("Received signal:", sig)
	}

	rpcServer.Terminate(ctx)
}

func setupAPIGateway() gateway.GatewayInterface {
	apiGateway := gateway.NewGateway(gatewayPort,
		constants.MaxHeaderBytes,
		constants.ReadTimeout,
		constants.WriteTimeout,
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   false,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
		runtime.WithErrorHandler(gateway.ExceptionHandler),
		runtime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, message proto.Message) error {
			md, ok := runtime.ServerMetadataFromContext(ctx)
			if !ok {
				return nil
			}

			if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
				code, err := strconv.Atoi(vals[0])
				if err != nil {
					return err
				}
				w.WriteHeader(code)
				delete(md.HeaderMD, "x-http-code")
				delete(w.Header(), "Grpc-Metadata-X-Http-Code")
			}

			return nil
		}),
	)

	return apiGateway
}

func Execute() {
	rootCmd.Flags().StringVarP(&rpcPort, "rpc", "r", "", "define rpc server port")
	rootCmd.Flags().StringVarP(&gatewayPort, "gateway", "g", "", "define gateway port")
	rootCmd.MarkFlagsRequiredTogether("rpc", "gateway")

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
