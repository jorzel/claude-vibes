package domain

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

// Executor is an interface that can be implemented by both *sql.DB and *sql.Tx
// This allows repositories to work with both direct database connections and transactions
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type EventRepository interface {
	Create(ctx context.Context, event *Event) error
	FindByID(ctx context.Context, id uuid.UUID) (*Event, error)
	FindAll(ctx context.Context) ([]*Event, error)
	Update(ctx context.Context, event *Event) error
	// Transaction-aware methods
	FindByIDWithLock(ctx context.Context, exec Executor, id uuid.UUID) (*Event, error)
	UpdateWithExecutor(ctx context.Context, exec Executor, event *Event) error
}

type BookingRepository interface {
	Create(ctx context.Context, booking *Booking) error
	FindByID(ctx context.Context, id uuid.UUID) (*Booking, error)
	// Transaction-aware methods
	CreateWithExecutor(ctx context.Context, exec Executor, booking *Booking) error
}
