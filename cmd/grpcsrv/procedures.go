package main

import (
	"context"
	"github.com/ChechenItza/booking/internal/booking"
	pb "github.com/ChechenItza/protobufs/gen/go/booking/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (bs *BookingServer) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.CreateBookingResponse, error) {
	if req.UserId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user_id must be positive")
	}
	if req.ResourceId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "resource_id must be positive")
	}
	if req.ResourceCapacity <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "resource_capacity must be positive")
	}
	if req.StartAt == nil || req.EndAt == nil {
		return nil, status.Errorf(codes.InvalidArgument, "start_at and end_at must be provided")
	}

	startAt := req.StartAt.AsTime()
	endAt := req.EndAt.AsTime()
	if !endAt.After(startAt) {
		return nil, status.Errorf(codes.InvalidArgument, "end_at must be after start_at")
	}

	id, err := bs.booking.Create(ctx, int(req.UserId), int(req.ResourceId), int(req.ResourceCapacity), startAt, endAt)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateBookingResponse{BookingId: int32(id)}, nil
}

func (bs *BookingServer) GetBookingsByResource(ctx context.Context, req *pb.GetBookingsByResourceRequest) (*pb.GetBookingsByResourceResponse, error) {
	if len(req.ResourceIds) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "resource_ids must be provided")
	}

	bookings, err := bs.booking.ListByResourceIds(ctx, req.ResourceIds)
	if err != nil {
		//TODO: switch on errors
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetBookingsByResourceResponse{Bookings: fromBookingInfoToGrpcInfo(bookings)}, nil
}

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
				Msg("failed to handle request")
		} else {
			bs.logger.Info().
				Str("method", info.FullMethod).
				Dur("duration", time.Since(start)).
				Msg("request completed")
		}
		return resp, err
	}
}

func fromBookingInfoToGrpcInfo(bookings []booking.Info) []*pb.BookingInfo {
	res := make([]*pb.BookingInfo, len(bookings))
	for i, b := range bookings {
		res[i] = &pb.BookingInfo{
			BookingId:  int32(b.Id),
			ResourceId: int32(b.ResourceId),
			StartAt:    timestamppb.New(b.StartAt),
			EndAt:      timestamppb.New(b.EndAt),
		}
	}

	return res
}
