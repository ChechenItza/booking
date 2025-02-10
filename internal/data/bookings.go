package data

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Booking struct {
	Id         int       `db:"id"`
	ResourceId int       `db:"resource_id"`
	StartAt    time.Time `db:"start_at"`
	EndAt      time.Time `db:"end_at"`
}

type BookingModel struct {
	pool *pgxpool.Pool
}

func (m *BookingModel) ListByResourceIds(ids []int) ([]Booking, error) {
	getQuery := `
		SELECT id, resource_id, start_at, end_at
		FROM bookings b
		WHERE b.resource_id in ($1)
	`
	rows, err := m.pool.Query(context.TODO(), getQuery, ids)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	bookings, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[Booking])
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (m *BookingModel) ListByUserId(userId int) ([]Booking, error) {
	getQuery := `
		SELECT id, resource_id, start_at, end_at
		FROM bookings b
		WHERE b.user_id = $1
	`
	rows, err := m.pool.Query(context.TODO(), getQuery, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	bookings, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[Booking])
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (m *BookingModel) Create(userId, resourceId, resourceCapacity int, startAt, endAt time.Time) (int, error) {
	tx, err := m.pool.BeginTx(context.TODO(), pgx.TxOptions{}) //TODO: no auto rollback on ctx timeout????
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(context.TODO()) //TODO: maybe log error

	// TODO: this is pessimistic locking, maybe think about changing to optimistic locking
	countQuery := `
		SELECT count
		FROM booking_count
		WHERE resource_id = $1
		FOR UPDATE
	`
	var count int
	err = tx.QueryRow(context.TODO(), countQuery, resourceId).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrRecordNotFound
		}
		return 0, err
	}
	if count >= resourceCapacity {
		return 0, ErrCapReached
	}

	insertQuery := `
		INSERT 
		INTO bookings (user_id, resource_id, start_at, end_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, resource_id, start_at, end_at) DO NOTHING
		RETURNING id
	`
	var id int
	err = tx.QueryRow(context.TODO(), insertQuery, userId, resourceId, startAt, endAt).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrAlreadyExists
		}
		return 0, err
	}

	incCountQuery := `
		UPDATE booking_count
		SET count = count + 1
		WHERE resource_id = $1 
	`
	ct, err := tx.Exec(context.TODO(), incCountQuery)
	if err != nil {
		return 0, err
	}
	if ct.RowsAffected() == 0 {
		return 0, ErrRecordNotFound
	}

	err = tx.Commit(context.TODO())
	if err != nil {
		return 0, err
	}

	return id, nil
}
