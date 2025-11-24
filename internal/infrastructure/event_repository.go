package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jorzel/booking-service/internal/domain"
)

type PostgresEventRepository struct {
	db DBClient
}

func NewPostgresEventRepository(db DBClient) *PostgresEventRepository {
	return &PostgresEventRepository{db: db}
}

func (r *PostgresEventRepository) Create(ctx context.Context, event *domain.Event) error {
	query := `
		INSERT INTO events (id, name, date, location, tickets)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		event.ID,
		event.Name,
		event.Date,
		event.Location,
		event.Tickets,
	)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

func (r *PostgresEventRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	query := `
		SELECT id, name, date, location, tickets
		FROM events
		WHERE id = $1
	`

	event := &domain.Event{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.Name,
		&event.Date,
		&event.Location,
		&event.Tickets,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrEventNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find event: %w", err)
	}

	return event, nil
}

func (r *PostgresEventRepository) FindAll(ctx context.Context) ([]*domain.Event, error) {
	query := `
		SELECT id, name, date, location, tickets
		FROM events
		ORDER BY date ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []*domain.Event
	for rows.Next() {
		event := &domain.Event{}
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Date,
			&event.Location,
			&event.Tickets,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %w", err)
	}

	return events, nil
}

func (r *PostgresEventRepository) Update(ctx context.Context, event *domain.Event) error {
	query := `
		UPDATE events
		SET name = $2, date = $3, location = $4, tickets = $5
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		event.ID,
		event.Name,
		event.Date,
		event.Location,
		event.Tickets,
	)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
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

// CreateWithExecutor creates an event using the provided executor (transaction or db)
func (r *PostgresEventRepository) CreateWithExecutor(ctx context.Context, exec domain.Executor, event *domain.Event) error {
	query := `
		INSERT INTO events (id, name, date, location, tickets)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := exec.ExecContext(
		ctx,
		query,
		event.ID,
		event.Name,
		event.Date,
		event.Location,
		event.Tickets,
	)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}
