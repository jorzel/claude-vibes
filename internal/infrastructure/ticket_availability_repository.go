package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jorzel/booking-service/internal/domain"
)

type PostgresTicketAvailabilityRepository struct {
	db DBClient
}

func NewPostgresTicketAvailabilityRepository(db DBClient) *PostgresTicketAvailabilityRepository {
	return &PostgresTicketAvailabilityRepository{db: db}
}

func (r *PostgresTicketAvailabilityRepository) Create(ctx context.Context, availability *domain.TicketAvailability) error {
	query := `
		INSERT INTO ticket_availability (event_id, available_tickets)
		VALUES ($1, $2)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		availability.EventID,
		availability.AvailableTickets,
	)
	if err != nil {
		return fmt.Errorf("failed to create ticket availability: %w", err)
	}

	return nil
}

func (r *PostgresTicketAvailabilityRepository) FindByEventID(ctx context.Context, eventID uuid.UUID) (*domain.TicketAvailability, error) {
	query := `
		SELECT event_id, available_tickets
		FROM ticket_availability
		WHERE event_id = $1
	`

	availability := &domain.TicketAvailability{}
	err := r.db.QueryRowContext(ctx, query, eventID).Scan(
		&availability.EventID,
		&availability.AvailableTickets,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrEventNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find ticket availability: %w", err)
	}

	return availability, nil
}

// CreateWithExecutor creates ticket availability using the provided executor (transaction or db)
func (r *PostgresTicketAvailabilityRepository) CreateWithExecutor(ctx context.Context, exec domain.Executor, availability *domain.TicketAvailability) error {
	query := `
		INSERT INTO ticket_availability (event_id, available_tickets)
		VALUES ($1, $2)
	`

	_, err := exec.ExecContext(
		ctx,
		query,
		availability.EventID,
		availability.AvailableTickets,
	)
	if err != nil {
		return fmt.Errorf("failed to create ticket availability: %w", err)
	}

	return nil
}

// FindByEventIDWithLock retrieves ticket availability by event ID with a row-level lock (FOR UPDATE)
// This should be used within a transaction to prevent concurrent modifications
func (r *PostgresTicketAvailabilityRepository) FindByEventIDWithLock(ctx context.Context, exec domain.Executor, eventID uuid.UUID) (*domain.TicketAvailability, error) {
	query := `
		SELECT event_id, available_tickets
		FROM ticket_availability
		WHERE event_id = $1
		FOR UPDATE
	`

	availability := &domain.TicketAvailability{}
	err := exec.QueryRowContext(ctx, query, eventID).Scan(
		&availability.EventID,
		&availability.AvailableTickets,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrEventNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find ticket availability: %w", err)
	}

	return availability, nil
}

// UpdateWithExecutor updates ticket availability using the provided executor (transaction or db)
func (r *PostgresTicketAvailabilityRepository) UpdateWithExecutor(ctx context.Context, exec domain.Executor, availability *domain.TicketAvailability) error {
	query := `
		UPDATE ticket_availability
		SET available_tickets = $2
		WHERE event_id = $1
	`

	result, err := exec.ExecContext(
		ctx,
		query,
		availability.EventID,
		availability.AvailableTickets,
	)
	if err != nil {
		return fmt.Errorf("failed to update ticket availability: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}
