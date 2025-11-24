package infrastructure

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/jorzel/booking-service/internal/domain"
)

// InstrumentedPostgresClient wraps sql.DB and tracks query metrics
type InstrumentedPostgresClient struct {
	*sql.DB
}

// NewInstrumentedPostgresClient creates a new instrumented postgres client
func NewInstrumentedPostgresClient(db *sql.DB) *InstrumentedPostgresClient {
	return &InstrumentedPostgresClient{DB: db}
}

// InstrumentedTx wraps sql.Tx and tracks query metrics
type InstrumentedTx struct {
	*sql.Tx
}

// ExecContext wraps the standard ExecContext with instrumentation
func (c *InstrumentedPostgresClient) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	operation := extractOperation(query)
	start := time.Now()

	result, err := c.DB.ExecContext(ctx, query, args...)

	duration := time.Since(start).Seconds()
	PostgresQueryDuration.WithLabelValues(operation).Observe(duration)

	status := "success"
	if err != nil {
		status = "error"
	}
	PostgresQueriesTotal.WithLabelValues(operation, status).Inc()

	return result, err
}

// QueryContext wraps the standard QueryContext with instrumentation
func (c *InstrumentedPostgresClient) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	operation := extractOperation(query)
	start := time.Now()

	rows, err := c.DB.QueryContext(ctx, query, args...)

	duration := time.Since(start).Seconds()
	PostgresQueryDuration.WithLabelValues(operation).Observe(duration)

	status := "success"
	if err != nil {
		status = "error"
	}
	PostgresQueriesTotal.WithLabelValues(operation, status).Inc()

	return rows, err
}

// QueryRowContext wraps the standard QueryRowContext with instrumentation
func (c *InstrumentedPostgresClient) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	operation := extractOperation(query)
	start := time.Now()

	row := c.DB.QueryRowContext(ctx, query, args...)

	duration := time.Since(start).Seconds()
	PostgresQueryDuration.WithLabelValues(operation).Observe(duration)
	PostgresQueriesTotal.WithLabelValues(operation, "success").Inc()

	return row
}

// BeginTx wraps the standard BeginTx and returns an instrumented transaction
func (c *InstrumentedPostgresClient) BeginTx(ctx context.Context, opts *sql.TxOptions) (domain.Transaction, error) {
	tx, err := c.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &InstrumentedTx{Tx: tx}, nil
}

// PingContext wraps the standard PingContext
func (c *InstrumentedPostgresClient) PingContext(ctx context.Context) error {
	return c.DB.PingContext(ctx)
}

// Close wraps the standard Close
func (c *InstrumentedPostgresClient) Close() error {
	return c.DB.Close()
}

// ExecContext wraps the transaction's ExecContext with instrumentation
func (tx *InstrumentedTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	operation := extractOperation(query)
	start := time.Now()

	result, err := tx.Tx.ExecContext(ctx, query, args...)

	duration := time.Since(start).Seconds()
	PostgresQueryDuration.WithLabelValues(operation).Observe(duration)

	status := "success"
	if err != nil {
		status = "error"
	}
	PostgresQueriesTotal.WithLabelValues(operation, status).Inc()

	return result, err
}

// QueryContext wraps the transaction's QueryContext with instrumentation
func (tx *InstrumentedTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	operation := extractOperation(query)
	start := time.Now()

	rows, err := tx.Tx.QueryContext(ctx, query, args...)

	duration := time.Since(start).Seconds()
	PostgresQueryDuration.WithLabelValues(operation).Observe(duration)

	status := "success"
	if err != nil {
		status = "error"
	}
	PostgresQueriesTotal.WithLabelValues(operation, status).Inc()

	return rows, err
}

// QueryRowContext wraps the transaction's QueryRowContext with instrumentation
func (tx *InstrumentedTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	operation := extractOperation(query)
	start := time.Now()

	row := tx.Tx.QueryRowContext(ctx, query, args...)

	duration := time.Since(start).Seconds()
	PostgresQueryDuration.WithLabelValues(operation).Observe(duration)
	PostgresQueriesTotal.WithLabelValues(operation, "success").Inc()

	return row
}

// extractOperation extracts the SQL operation type from a query string
func extractOperation(query string) string {
	// Trim whitespace and convert to uppercase
	trimmed := strings.TrimSpace(query)
	upper := strings.ToUpper(trimmed)

	// Extract the first word (operation)
	parts := strings.Fields(upper)
	if len(parts) == 0 {
		return "UNKNOWN"
	}

	operation := parts[0]

	// Normalize common operations
	switch operation {
	case "SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "DROP", "ALTER", "TRUNCATE":
		return operation
	case "WITH":
		// For CTEs, find the actual operation
		for _, part := range parts[1:] {
			if part == "SELECT" || part == "INSERT" || part == "UPDATE" || part == "DELETE" {
				return part
			}
		}
		return "WITH"
	default:
		return "OTHER"
	}
}
