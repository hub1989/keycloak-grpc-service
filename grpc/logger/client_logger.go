package logger

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func ClientLogger(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Info("outgoing gRpc request >> ", method)
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		log.Printf("method %s failed: %s", method, err)
	}
	return err
}
