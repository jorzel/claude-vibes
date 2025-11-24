package tests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jorzel/booking-service/internal/app"
	"github.com/jorzel/booking-service/internal/domain"
	"github.com/jorzel/booking-service/internal/infrastructure"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(60 * time.Second),
	}

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := postgres.Host(ctx)
	require.NoError(t, err)

	port, err := postgres.MappedPort(ctx, "5432")
	require.NoError(t, err)

	config := infrastructure.Config{
		Host:     host,
		Port:     port.Int(),
		User:     "test",
		Password: "test",
		Database: "testdb",
		SSLMode:  "disable",
	}

	db, err := infrastructure.NewPostgresDB(config)
	require.NoError(t, err)

	migrationSQL, err := os.ReadFile("../internal/infrastructure/migrations/001_create_tables.sql")
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, string(migrationSQL))
	require.NoError(t, err)

	cleanup := func() {
		db.Close()
		postgres.Terminate(ctx)
	}

	return db, cleanup
}

func TestEventService_Integration(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	eventRepo := infrastructure.NewPostgresEventRepository(db)
	eventService := app.NewEventService(eventRepo, logger)

	ctx := context.Background()

	t.Run("creates and retrieves event", func(t *testing.T) {
		req := app.CreateEventRequest{
			Name:     "Summer Festival",
			Date:     time.Now().Add(30 * 24 * time.Hour),
			Location: "Central Park",
			Tickets:  200,
		}

		created, err := eventService.CreateEvent(ctx, req)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, created.ID)
		assert.Equal(t, req.Name, created.Name)
		assert.Equal(t, req.Location, created.Location)
		assert.Equal(t, req.Tickets, created.Tickets)

		retrieved, err := eventService.GetEvent(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.Name, retrieved.Name)
	})

	t.Run("lists all events", func(t *testing.T) {
		events, err := eventService.ListEvents(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, events)
	})

	t.Run("returns error for non-existent event", func(t *testing.T) {
		nonExistentID := uuid.New()
		_, err := eventService.GetEvent(ctx, nonExistentID)
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrEventNotFound)
	})
}

func TestBookingService_Integration(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	eventRepo := infrastructure.NewPostgresEventRepository(db)
	bookingRepo := infrastructure.NewPostgresBookingRepository(db)
	eventService := app.NewEventService(eventRepo, logger)
	bookingService := app.NewBookingService(bookingRepo, eventRepo, db, logger)

	ctx := context.Background()

	t.Run("creates booking and decrements available tickets", func(t *testing.T) {
		event, err := eventService.CreateEvent(ctx, app.CreateEventRequest{
			Name:     "Rock Concert",
			Date:     time.Now().Add(15 * 24 * time.Hour),
			Location: "Stadium",
			Tickets:  100,
		})
		require.NoError(t, err)

		userID := uuid.New()
		bookingReq := app.CreateBookingRequest{
			EventID:       event.ID,
			UserID:        userID,
			TicketsBooked: 5,
		}

		booking, err := bookingService.CreateBooking(ctx, bookingReq)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, booking.ID)
		assert.Equal(t, event.ID, booking.EventID)
		assert.Equal(t, userID, booking.UserID)
		assert.Equal(t, 5, booking.TicketsBooked)

		updatedEvent, err := eventService.GetEvent(ctx, event.ID)
		require.NoError(t, err)
		assert.Equal(t, 95, updatedEvent.AvailableTickets)
	})

	t.Run("returns error when booking more tickets than available", func(t *testing.T) {
		event, err := eventService.CreateEvent(ctx, app.CreateEventRequest{
			Name:     "Small Event",
			Date:     time.Now().Add(7 * 24 * time.Hour),
			Location: "Small Venue",
			Tickets:  5,
		})
		require.NoError(t, err)

		bookingReq := app.CreateBookingRequest{
			EventID:       event.ID,
			UserID:        uuid.New(),
			TicketsBooked: 10,
		}

		_, err = bookingService.CreateBooking(ctx, bookingReq)
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInsufficientTickets)
	})

	t.Run("retrieves booking by id", func(t *testing.T) {
		event, err := eventService.CreateEvent(ctx, app.CreateEventRequest{
			Name:     "Theater Show",
			Date:     time.Now().Add(10 * 24 * time.Hour),
			Location: "Theater",
			Tickets:  50,
		})
		require.NoError(t, err)

		created, err := bookingService.CreateBooking(ctx, app.CreateBookingRequest{
			EventID:       event.ID,
			UserID:        uuid.New(),
			TicketsBooked: 2,
		})
		require.NoError(t, err)

		retrieved, err := bookingService.GetBooking(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.EventID, retrieved.EventID)
		assert.Equal(t, created.TicketsBooked, retrieved.TicketsBooked)
	})

	t.Run("handles concurrent bookings correctly", func(t *testing.T) {
		event, err := eventService.CreateEvent(ctx, app.CreateEventRequest{
			Name:     "Popular Concert",
			Date:     time.Now().Add(20 * 24 * time.Hour),
			Location: "Arena",
			Tickets:  10,
		})
		require.NoError(t, err)

		errChan := make(chan error, 3)
		for i := 0; i < 3; i++ {
			go func() {
				_, err := bookingService.CreateBooking(ctx, app.CreateBookingRequest{
					EventID:       event.ID,
					UserID:        uuid.New(),
					TicketsBooked: 4,
				})
				errChan <- err
			}()
		}

		successCount := 0
		failureCount := 0
		for i := 0; i < 3; i++ {
			err := <-errChan
			if err == nil {
				successCount++
			} else {
				failureCount++
			}
		}

		assert.Equal(t, 2, successCount, "expected 2 successful bookings")
		assert.Equal(t, 1, failureCount, "expected 1 failed booking")

		finalEvent, err := eventService.GetEvent(ctx, event.ID)
		require.NoError(t, err)
		assert.Equal(t, 2, finalEvent.AvailableTickets)
	})
}

func TestHTTPEndpoints_Integration(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	eventRepo := infrastructure.NewPostgresEventRepository(db)
	bookingRepo := infrastructure.NewPostgresBookingRepository(db)
	eventService := app.NewEventService(eventRepo, logger)
	bookingService := app.NewBookingService(bookingRepo, eventRepo, db, logger)

	ctx := context.Background()

	t.Run("full booking flow", func(t *testing.T) {
		event, err := eventService.CreateEvent(ctx, app.CreateEventRequest{
			Name:     "Integration Test Event",
			Date:     time.Now().Add(5 * 24 * time.Hour),
			Location: "Test Location",
			Tickets:  50,
		})
		require.NoError(t, err)

		events, err := eventService.ListEvents(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, events)

		var found bool
		for _, e := range events {
			if e.ID == event.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "created event should appear in list")

		booking, err := bookingService.CreateBooking(ctx, app.CreateBookingRequest{
			EventID:       event.ID,
			UserID:        uuid.New(),
			TicketsBooked: 3,
		})
		require.NoError(t, err)

		retrievedBooking, err := bookingService.GetBooking(ctx, booking.ID)
		require.NoError(t, err)
		assert.Equal(t, booking.ID, retrievedBooking.ID)

		updatedEvent, err := eventService.GetEvent(ctx, event.ID)
		require.NoError(t, err)
		assert.Equal(t, 47, updatedEvent.AvailableTickets)
	})
}

func BenchmarkCreateBooking(b *testing.B) {
	db, cleanup := setupBenchDB(b)
	defer cleanup()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	eventRepo := infrastructure.NewPostgresEventRepository(db)
	bookingRepo := infrastructure.NewPostgresBookingRepository(db)
	eventService := app.NewEventService(eventRepo, logger)
	bookingService := app.NewBookingService(bookingRepo, eventRepo, db, logger)

	ctx := context.Background()

	event, err := eventService.CreateEvent(ctx, app.CreateEventRequest{
		Name:     "Benchmark Event",
		Date:     time.Now().Add(30 * 24 * time.Hour),
		Location: "Benchmark Location",
		Tickets:  10000,
	})
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := bookingService.CreateBooking(ctx, app.CreateBookingRequest{
			EventID:       event.ID,
			UserID:        uuid.New(),
			TicketsBooked: 1,
		})
		if err != nil {
			b.Fatalf("booking failed: %v", err)
		}
	}
}

func setupBenchDB(b *testing.B) (*sql.DB, func()) {
	b.Helper()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "benchdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(60 * time.Second),
	}

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(b, err)

	host, err := postgres.Host(ctx)
	require.NoError(b, err)

	port, err := postgres.MappedPort(ctx, "5432")
	require.NoError(b, err)

	dsn := fmt.Sprintf("host=%s port=%d user=test password=test dbname=benchdb sslmode=disable",
		host, port.Int())
	db, err := sql.Open("postgres", dsn)
	require.NoError(b, err)

	migrationSQL, err := os.ReadFile("../internal/infrastructure/migrations/001_create_tables.sql")
	require.NoError(b, err)

	_, err = db.ExecContext(ctx, string(migrationSQL))
	require.NoError(b, err)

	cleanup := func() {
		db.Close()
		postgres.Terminate(ctx)
	}

	return db, cleanup
}
