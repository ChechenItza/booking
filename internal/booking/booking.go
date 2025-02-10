package booking

import (
	"github.com/ChechenItza/resource-booking/internal/data"
	"time"
)

type BookingInfo struct {
	Id         int
	ResourceId int
	StartAt    time.Time
	EndAt      time.Time
}

//type BookingService interface {
//	Create(userId int, resourceId int, startAt time.Time, endAt time.Time) (int, error)
//	ListByUserId(userId int) ([]BookingInfo, error)
//	ListByResourceIds(resourceIds []int) ([]BookingInfo, error)
//	Remove(bookingId int) error
//}

/*
QUESTIONS:
1. transactions between repo functions
*/

type BookingService struct {
	models data.Models
}

func NewBookingService(models data.Models) BookingService {
	return BookingService{models}
}

func (b *BookingService) Create(userId, resourceId, resourceCapacity int, startAt, endAt time.Time) (int, error) {
	id, err := b.models.Bookings.Create(userId, resourceId, resourceCapacity, startAt, endAt)
	if err != nil {
		//TODO: retry in certain conditions
		return 0, err
	}
	
	return id, nil
}

func (b *BookingService) ListByUserId(userId int) ([]BookingInfo, error) {
	bookings, err := b.models.Bookings.ListByUserId(userId)
	if err != nil {
		return nil, err
	}

	return fromDataBookingsToBookingInfos(bookings), nil
}

func (b *BookingService) ListByResourceIds(resourceIds []int) ([]BookingInfo, error) {
	bookings, err := b.models.Bookings.ListByResourceIds(resourceIds)
	if err != nil {
		return nil, err
	}

	return fromDataBookingsToBookingInfos(bookings), nil
}

func fromDataBookingToBookingInfo(booking data.Booking) BookingInfo {
	return BookingInfo{
		Id:         booking.Id,
		ResourceId: booking.ResourceId,
		StartAt:    booking.StartAt,
		EndAt:      booking.EndAt,
	}
}

func fromDataBookingsToBookingInfos(bookings []data.Booking) []BookingInfo {
	bookingInfos := make([]BookingInfo, len(bookings))
	for i, booking := range bookings {
		bookingInfos[i] = fromDataBookingToBookingInfo(booking)
	}
	return bookingInfos
}
