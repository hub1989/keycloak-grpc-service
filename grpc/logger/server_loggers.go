package logger

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func ServerLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.WithField("method", info.FullMethod).Info("incoming gRpc request")
	resp, err := handler(ctx, req)
	if err != nil {
		log.WithFields(log.Fields{
			"method": info.FullMethod,
		}).WithError(err).Error("method failed")
	}
	return resp, err
}
