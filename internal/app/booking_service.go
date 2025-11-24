package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jorzel/booking-service/internal/domain"
	"github.com/rs/zerolog"
)

type BookingService struct {
	bookingRepo domain.BookingRepository
	eventRepo   domain.EventRepository
	db          *sql.DB
	logger      zerolog.Logger
}

func NewBookingService(
	bookingRepo domain.BookingRepository,
	eventRepo domain.EventRepository,
	db *sql.DB,
	logger zerolog.Logger,
) *BookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		eventRepo:   eventRepo,
		db:          db,
		logger:      logger.With().Str("service", "booking").Logger(),
	}
}

type CreateBookingRequest struct {
	EventID       uuid.UUID
	UserID        uuid.UUID
	TicketsBooked int
}

func (s *BookingService) CreateBooking(ctx context.Context, req CreateBookingRequest) (*domain.Booking, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to begin transaction")
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	event, err := s.eventRepo.FindByID(ctx, req.EventID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("event_id", req.EventID.String()).
			Msg("failed to find event")
		return nil, fmt.Errorf("failed to find event: %w", err)
	}

	if err := event.ReserveTickets(req.TicketsBooked); err != nil {
		s.logger.Warn().
			Err(err).
			Str("event_id", req.EventID.String()).
			Int("requested", req.TicketsBooked).
			Int("available", event.AvailableTickets).
			Msg("insufficient tickets")
		return nil, err
	}

	if err := s.eventRepo.Update(ctx, event); err != nil {
		s.logger.Error().
			Err(err).
			Str("event_id", event.ID.String()).
			Msg("failed to update event")
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	booking, err := domain.NewBooking(req.EventID, req.UserID, req.TicketsBooked)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create booking domain object")
		return nil, fmt.Errorf("invalid booking data: %w", err)
	}

	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		s.logger.Error().
			Err(err).
			Str("booking_id", booking.ID.String()).
			Msg("failed to save booking")
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error().Err(err).Msg("failed to commit transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info().
		Str("booking_id", booking.ID.String()).
		Str("event_id", booking.EventID.String()).
		Str("user_id", booking.UserID.String()).
		Int("tickets", booking.TicketsBooked).
		Msg("booking created")

	return booking, nil
}

func (s *BookingService) GetBooking(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	booking, err := s.bookingRepo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Str("booking_id", id.String()).Msg("failed to find booking")
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	return booking, nil
}
