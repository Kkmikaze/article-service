package rpcclient

import (
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RPCClient struct {
	UnaryConn *grpc.ClientConn
}

func NewRPCClient(
	addr string,
	unaryInterceptors []grpc.UnaryClientInterceptor,
	opts ...grpc.DialOption,
) (*RPCClient, error) {
	rpc := &RPCClient{}

	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	opts = append(opts, creds)

	var chainUnaryInterceptors []grpc.DialOption
	chainUnaryInterceptors = append(chainUnaryInterceptors, opts...)
	if len(unaryInterceptors) > 0 {
		chainUnaryInterceptors = append(chainUnaryInterceptors, grpc.WithChainUnaryInterceptor(unaryInterceptors...))
	}

	var host string
	if strings.Contains(addr, ":") {
		host = addr
	} else {
		host = fmt.Sprintf(":%v", addr)
	}

	unaryConn, errUnary := grpc.NewClient(host, chainUnaryInterceptors...)
	if errUnary != nil {
		return nil, errUnary
	}

	rpc.UnaryConn = unaryConn

	return rpc, nil
}