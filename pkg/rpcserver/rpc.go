package rpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type RPCServer struct {
	*grpc.Server
	Addr     string
	Listener net.Listener
	Network  string
}

type RPCServerInterface interface {
	Run() error
	StopListener()
	Terminate(ctx context.Context)
}

func (rpc *RPCServer) Run() error {
	var err error

	var host string
	if strings.Contains(rpc.Addr, ":") {
		host = rpc.Addr
	} else {
		host = fmt.Sprintf(":%v", rpc.Addr)
	}

	rpc.Listener, err = net.Listen(rpc.Network, host)
	if err != nil {
		return err
	}

	go func() {
		if err := rpc.Serve(rpc.Listener); err != nil {
			log.Fatalf("Server exited with error: %v\n", err)
		}
	}()

	return nil
}

func (rpc *RPCServer) StopListener() {
	if err := rpc.Listener.Close(); err != nil {
		log.Fatalf("Failed to close %s %s: %v", rpc.Network, rpc.Addr, err)
	}
}

func (rpc *RPCServer) Terminate(ctx context.Context) {
	go func() {
		defer rpc.Server.GracefulStop()
		<-ctx.Done()
	}()
}

func NewRPCServer(addr,
	network string,
	tls bool,
	unaryInterceptors []grpc.UnaryServerInterceptor,
	streamInterceptors []grpc.StreamServerInterceptor,
	opts ...grpc.ServerOption,
) *RPCServer {
	if tls {
		certFile := "ssl/certificates/server.crt" // => your certFile file path
		keyFile := "ssl/server.pem"               // => your keyFile file path

		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		opts = append(opts, grpc.Creds(creds))
	}

	if len(unaryInterceptors) > 0 {
		opts = append(opts, grpc.ChainUnaryInterceptor(unaryInterceptors...))
	}

	if len(streamInterceptors) > 0 {
		opts = append(opts, grpc.ChainStreamInterceptor(streamInterceptors...))
	}

	s := grpc.NewServer(opts...)

	return &RPCServer{
		Addr:    addr,
		Network: network,
		Server:  s,
	}
}