package booking

import (
	"context"
	"github.com/ChechenItza/booking/internal/data"
	"time"
)

type Info struct {
	Id         int
	ResourceId int
	StartAt    time.Time
	EndAt      time.Time
}

type Service struct {
	models data.Models
}

func NewService(models data.Models) Service {
	return Service{models}
}

func (b *Service) Create(ctx context.Context, userId, resourceId, resourceCapacity int, startAt, endAt time.Time) (int, error) {
	id, err := b.models.Bookings.Create(ctx, userId, resourceId, resourceCapacity, startAt, endAt)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (b *Service) ListByResourceIds(ctx context.Context, resourceIds []int32) ([]Info, error) {
	bookings, err := b.models.Bookings.ListByResourceIds(ctx, resourceIds)
	if err != nil {
		return nil, err
	}

	return fromDataBookingsToBookingInfos(bookings), nil
}

func fromDataBookingToBookingInfo(booking data.Booking) Info {
	return Info{
		Id:         booking.Id,
		ResourceId: booking.ResourceId,
		StartAt:    booking.StartAt,
		EndAt:      booking.EndAt,
	}
}

func fromDataBookingsToBookingInfos(bookings []data.Booking) []Info {
	bookingInfos := make([]Info, len(bookings))
	for i, booking := range bookings {
		bookingInfos[i] = fromDataBookingToBookingInfo(booking)
	}
	return bookingInfos
}
