package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jorzel/booking-service/internal/domain"
	"github.com/jorzel/booking-service/internal/infrastructure"
	"github.com/rs/zerolog"
)

type EventService struct {
	repo                   domain.EventRepository
	ticketAvailabilityRepo domain.TicketAvailabilityRepository
	db                     infrastructure.DBClient
	logger                 zerolog.Logger
}

func NewEventService(
	repo domain.EventRepository,
	ticketAvailabilityRepo domain.TicketAvailabilityRepository,
	db infrastructure.DBClient,
	logger zerolog.Logger,
) *EventService {
	return &EventService{
		repo:                   repo,
		ticketAvailabilityRepo: ticketAvailabilityRepo,
		db:                     db,
		logger:                 logger.With().Str("service", "event").Logger(),
	}
}

type CreateEventRequest struct {
	Name     string
	Date     time.Time
	Location string
	Tickets  int
}

func (s *EventService) CreateEvent(ctx context.Context, req CreateEventRequest) (*domain.Event, error) {
	event, err := domain.NewEvent(req.Name, req.Location, req.Date, req.Tickets)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create event domain object")
		return nil, fmt.Errorf("invalid event data: %w", err)
	}

	// Create TicketAvailability aggregate for the event
	ticketAvailability, err := domain.NewTicketAvailability(event.ID, req.Tickets)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create ticket availability domain object")
		return nil, fmt.Errorf("invalid ticket availability data: %w", err)
	}

	// Use transaction to ensure atomic creation of both Event and TicketAvailability
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to begin transaction")
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := s.repo.CreateWithExecutor(ctx, tx, event); err != nil {
		s.logger.Error().Err(err).Str("event_id", event.ID.String()).Msg("failed to save event")
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	if err := s.ticketAvailabilityRepo.CreateWithExecutor(ctx, tx, ticketAvailability); err != nil {
		s.logger.Error().Err(err).Str("event_id", event.ID.String()).Msg("failed to save ticket availability")
		return nil, fmt.Errorf("failed to create ticket availability: %w", err)
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error().Err(err).Msg("failed to commit transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info().
		Str("event_id", event.ID.String()).
		Str("name", event.Name).
		Int("tickets", event.Tickets).
		Msg("event and ticket availability created")

	return event, nil
}

func (s *EventService) GetEvent(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	event, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Str("event_id", id.String()).Msg("failed to find event")
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return event, nil
}

func (s *EventService) ListEvents(ctx context.Context) ([]*domain.Event, error) {
	events, err := s.repo.FindAll(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to list events")
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	s.logger.Debug().Int("count", len(events)).Msg("events listed")
	return events, nil
}
