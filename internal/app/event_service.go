package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jorzel/booking-service/internal/domain"
	"github.com/rs/zerolog"
)

type EventService struct {
	repo   domain.EventRepository
	logger zerolog.Logger
}

func NewEventService(repo domain.EventRepository, logger zerolog.Logger) *EventService {
	return &EventService{
		repo:   repo,
		logger: logger.With().Str("service", "event").Logger(),
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

	if err := s.repo.Create(ctx, event); err != nil {
		s.logger.Error().Err(err).Str("event_id", event.ID.String()).Msg("failed to save event")
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	s.logger.Info().
		Str("event_id", event.ID.String()).
		Str("name", event.Name).
		Int("tickets", event.Tickets).
		Msg("event created")

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
