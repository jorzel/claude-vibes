package infrastructure

import (
	"context"
	"database/sql"

	"github.com/jorzel/booking-service/internal/domain"
)

// DBClient defines the interface for database operations
// This interface is implemented by both raw sql.DB (via DBClientAdapter) and InstrumentedPostgresClient
type DBClient interface {
	// ExecContext executes a query without returning any rows
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// QueryContext executes a query that returns rows
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	// QueryRowContext executes a query that returns at most one row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row

	// BeginTx starts a transaction
	BeginTx(ctx context.Context, opts *sql.TxOptions) (domain.Transaction, error)

	// PingContext verifies a connection to the database
	PingContext(ctx context.Context) error

	// Close closes the database connection
	Close() error
}

// DBClientAdapter wraps sql.DB to implement the DBClient interface
// This allows using raw sql.DB where DBClient is expected (useful for testing or non-instrumented scenarios)
type DBClientAdapter struct {
	db *sql.DB
}

// NewDBClientAdapter creates a new adapter for sql.DB
func NewDBClientAdapter(db *sql.DB) *DBClientAdapter {
	return &DBClientAdapter{db: db}
}

func (a *DBClientAdapter) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return a.db.ExecContext(ctx, query, args...)
}

func (a *DBClientAdapter) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return a.db.QueryContext(ctx, query, args...)
}

func (a *DBClientAdapter) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return a.db.QueryRowContext(ctx, query, args...)
}

func (a *DBClientAdapter) BeginTx(ctx context.Context, opts *sql.TxOptions) (domain.Transaction, error) {
	return a.db.BeginTx(ctx, opts)
}

func (a *DBClientAdapter) PingContext(ctx context.Context) error {
	return a.db.PingContext(ctx)
}

func (a *DBClientAdapter) Close() error {
	return a.db.Close()
}
