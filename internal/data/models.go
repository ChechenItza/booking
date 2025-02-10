package data

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrCapReached     = errors.New("capacity reached")
	ErrAlreadyExists  = errors.New("already exists")
	ErrTimeConflict   = errors.New("conflicting time already booked on the resource")
)

type Models struct {
	Bookings BookingModel
}

func NewModels(pool *pgxpool.Pool) Models {
	return Models{
		Bookings: BookingModel{pool},
	}
}
