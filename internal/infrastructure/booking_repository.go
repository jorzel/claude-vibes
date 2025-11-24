package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jorzel/booking-service/internal/domain"
)

type PostgresBookingRepository struct {
	db *sql.DB
}

func NewPostgresBookingRepository(db *sql.DB) *PostgresBookingRepository {
	return &PostgresBookingRepository{db: db}
}

func (r *PostgresBookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	query := `
		INSERT INTO bookings (id, event_id, user_id, tickets_booked, booked_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		booking.ID,
		booking.EventID,
		booking.UserID,
		booking.TicketsBooked,
		booking.BookedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}

	return nil
}

func (r *PostgresBookingRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	query := `
		SELECT id, event_id, user_id, tickets_booked, booked_at
		FROM bookings
		WHERE id = $1
	`

	booking := &domain.Booking{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&booking.ID,
		&booking.EventID,
		&booking.UserID,
		&booking.TicketsBooked,
		&booking.BookedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrBookingNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find booking: %w", err)
	}

	return booking, nil
}

// CreateWithExecutor creates a booking using the provided executor (transaction or db)
func (r *PostgresBookingRepository) CreateWithExecutor(ctx context.Context, exec domain.Executor, booking *domain.Booking) error {
	query := `
		INSERT INTO bookings (id, event_id, user_id, tickets_booked, booked_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := exec.ExecContext(
		ctx,
		query,
		booking.ID,
		booking.EventID,
		booking.UserID,
		booking.TicketsBooked,
		booking.BookedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}

	return nil
}
