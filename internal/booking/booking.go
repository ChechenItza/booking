package booking

import (
	"context"
	"errors"
	"github.com/ChechenItza/booking/internal/data"
	"time"
)

var (
	ErrStorageTimeout = errors.New("storage timeout")
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
	type CreateResult struct {
		id  int
		err error
	}

	res := make(chan CreateResult)
	go func() {
		id, err := b.models.Bookings.Create(ctx, userId, resourceId, resourceCapacity, startAt, endAt)
		res <- CreateResult{id, err}
	}()

	select {
	case <-ctx.Done():
		return 0, ErrStorageTimeout
	case r := <-res:
		return r.id, r.err
	}
}

func (b *Service) ListByUserId(userId int) ([]Info, error) {
	bookings, err := b.models.Bookings.ListByUserId(userId)
	if err != nil {
		return nil, err
	}

	return fromDataBookingsToBookingInfos(bookings), nil
}

func (b *Service) ListByResourceIds(ctx context.Context, resourceIds []int32) ([]Info, error) {
	type ListByResourceIdsRes struct {
		bookings []data.Booking
		err      error
	}

	res := make(chan ListByResourceIdsRes)
	go func() {
		bookings, err := b.models.Bookings.ListByResourceIds(resourceIds)
		res <- ListByResourceIdsRes{bookings, err}
	}()

	select {
	case <-ctx.Done():
		return nil, ErrStorageTimeout
	case r := <-res:
		return fromDataBookingsToBookingInfos(r.bookings), nil
	}
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
