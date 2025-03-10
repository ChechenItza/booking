package main

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

func (bs *BookingServer) LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		bs.logger.Info().Str("method", info.FullMethod).Msg("incoming request")

		resp, err := handler(ctx, req)

		if err != nil {
			bs.logger.Error().
				Str("method", info.FullMethod).
				Err(err).
				Dur("duration", time.Since(start)).
				Msg("request failed")
		} else {
			bs.logger.Info().
				Str("method", info.FullMethod).
				Dur("duration", time.Since(start)).
				Msg("request completed")
		}
		return resp, err
	}
}
