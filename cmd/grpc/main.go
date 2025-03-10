package main

import (
	"context"
	"github.com/ChechenItza/booking/internal/booking"
	"github.com/ChechenItza/booking/internal/data"
	pb "github.com/ChechenItza/protobufs/gen/go/booking/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type BookingServer struct {
	pb.UnimplementedBookingServiceServer
	booking booking.Service
	logger  zerolog.Logger
}

func main() {
	prettyLogger := zerolog.NewConsoleWriter()
	logger := zerolog.New(prettyLogger).With().Timestamp().Logger()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to listen")
	}

	pool, err := openDB("postgresql://admin:admin@127.0.0.1:5432/booking")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to db")
	}

	bookingService := booking.NewService(data.NewModels(pool))

	srv := BookingServer{
		logger:  logger,
		booking: bookingService,
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(srv.LoggingInterceptor()),
	)

	pb.RegisterBookingServiceServer(grpcServer, &srv)
	srv.logger.Info().Msg("gRPC server is running on :50051")

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal().Err(err).Msg("failed to serve")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Warn().Msg("Shutdown signal received, initiating graceful shutdown")
	grpcServer.GracefulStop()

	logger.Warn().Msg("gRPC server shutdown complete")
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	dbCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	//dbCfg.MaxConns = int32(cfg.Db.MaxConns)
	//dbCfg.MinConns = int32(cfg.Db.MinConns)
	//
	//duration, err := time.ParseDuration(cfg.Db.MaxIdleTime)
	//if err != nil {
	//	return nil, err
	//}
	//dbCfg.MaxConnIdleTime = duration

	pool, err := pgxpool.NewWithConfig(context.Background(), dbCfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
