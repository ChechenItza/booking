package data

import (
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (m *BookingModel) ListByResourceIds(ctx context.Context, ids []int32) ([]Booking, error) {
	getQuery := `
		SELECT id, resource_id, start_at, end_at
		FROM bookings b
		WHERE b.resource_id = any($1)
	`
	rows, err := m.pool.Query(ctx, getQuery, ids)
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

func (m *BookingModel) ListByUserId(ctx context.Context, userId int) ([]Booking, error) {
	getQuery := `
		SELECT id, resource_id, start_at, end_at
		FROM bookings b
		WHERE b.user_id = $1
	`
	rows, err := m.pool.Query(ctx, getQuery, userId)
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

func (m *BookingModel) Create(ctx context.Context, userId, resourceId, cap int, startAt, endAt time.Time) (int, error) {
	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{}) //TODO: no auto rollback on ctx timeout????
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx) //TODO: maybe log error

	// TODO: this is pessimistic locking, maybe think about changing to optimistic locking
	countQuery := `
		SELECT count
		FROM booking_count
		WHERE resource_id = $1
		FOR UPDATE
	`
	var count int
	err = tx.QueryRow(ctx, countQuery, resourceId).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrRecordNotFound
		}
		return 0, err
	}
	if count >= cap {
		return 0, ErrCapReached
	}

	insertQuery := `
		INSERT 
		INTO bookings (user_id, resource_id, start_at, end_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int
	err = tx.QueryRow(ctx, insertQuery, userId, resourceId, startAt, endAt).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.ExclusionViolation {
				return 0, ErrTimeConflict
			}
		}
		return 0, err
	}

	incCountQuery := `
		UPDATE booking_count
		SET count = count + 1, updated_at = now()
		WHERE resource_id = $1 
	`
	_, err = tx.Exec(ctx, incCountQuery, resourceId)
	if err != nil {
		return 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return id, nil
}
