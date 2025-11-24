package tests

import (
	"context"
	"testing"
	"time"

	"github.com/jorzel/booking-service/internal/domain"
	"github.com/jorzel/booking-service/internal/infrastructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEventRepository_CreateEvent_WithAvailableTicketsColumn tests the scenario where
// migration 002 didn't run, leaving the available_tickets column in the events table
// This replicates the production bug
func TestEventRepository_CreateEvent_WithAvailableTicketsColumn(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Re-add the available_tickets column to simulate the production state
	// where migration 002 didn't complete
	_, err := db.ExecContext(ctx, `
		ALTER TABLE events ADD COLUMN IF NOT EXISTS available_tickets INT NOT NULL DEFAULT 0;
	`)
	require.NoError(t, err)

	dbClient := infrastructure.NewDBClientAdapter(db)
	eventRepo := infrastructure.NewPostgresEventRepository(dbClient)

	event, err := domain.NewEvent("Test Event", "Test Location", time.Now().Add(24*time.Hour), 100)
	require.NoError(t, err)

	// With the old code (without available_tickets in INSERT), this would succeed
	// but use the DEFAULT value of 0, which is incorrect
	err = eventRepo.Create(ctx, event)
	require.NoError(t, err)

	// Verify the bug: available_tickets is 0 instead of 100
	var availableTickets int
	err = db.QueryRowContext(ctx, `
		SELECT available_tickets
		FROM events
		WHERE id = $1
	`, event.ID).Scan(&availableTickets)

	require.NoError(t, err)
	// This demonstrates the bug - available_tickets should be 100 but is 0
	assert.Equal(t, 0, availableTickets, "This test demonstrates the bug: available_tickets defaults to 0")
}

// TestEventRepository_AfterProperMigrations verifies that after running all migrations correctly,
// the available_tickets column is removed and events can be created successfully
func TestEventRepository_AfterProperMigrations(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Verify that available_tickets column doesn't exist after migration 002
	var columnExists bool
	err := db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.columns
			WHERE table_name='events' AND column_name='available_tickets'
		)
	`).Scan(&columnExists)
	require.NoError(t, err)
	assert.False(t, columnExists, "available_tickets column should not exist after migration 002")

	// Verify that ticket_availability table exists
	var tableExists bool
	err = db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_name='ticket_availability'
		)
	`).Scan(&tableExists)
	require.NoError(t, err)
	assert.True(t, tableExists, "ticket_availability table should exist after migration 002")

	// Now create an event and it should work correctly
	dbClient := infrastructure.NewDBClientAdapter(db)
	eventRepo := infrastructure.NewPostgresEventRepository(dbClient)

	event, err := domain.NewEvent("Test Event", "Test Location", time.Now().Add(24*time.Hour), 100)
	require.NoError(t, err)

	err = eventRepo.Create(ctx, event)
	require.NoError(t, err, "Creating event should succeed after proper migrations")

	// Verify the event was created correctly
	retrieved, err := eventRepo.FindByID(ctx, event.ID)
	require.NoError(t, err)
	assert.Equal(t, event.ID, retrieved.ID)
	assert.Equal(t, event.Tickets, retrieved.Tickets)
}
