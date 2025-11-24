package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jorzel/booking-service/internal/domain"
	"github.com/jorzel/booking-service/internal/infrastructure"
	"github.com/rs/zerolog"
)

type BookingService struct {
	bookingRepo            domain.BookingRepository
	ticketAvailabilityRepo domain.TicketAvailabilityRepository
	db                     infrastructure.DBClient
	logger                 zerolog.Logger
}

func NewBookingService(
	bookingRepo domain.BookingRepository,
	ticketAvailabilityRepo domain.TicketAvailabilityRepository,
	db infrastructure.DBClient,
	logger zerolog.Logger,
) *BookingService {
	return &BookingService{
		bookingRepo:            bookingRepo,
		ticketAvailabilityRepo: ticketAvailabilityRepo,
		db:                     db,
		logger:                 logger.With().Str("service", "booking").Logger(),
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

	// Lock the TicketAvailability aggregate (not the Event entity)
	ticketAvailability, err := s.ticketAvailabilityRepo.FindByEventIDWithLock(ctx, tx, req.EventID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("event_id", req.EventID.String()).
			Msg("failed to find ticket availability")
		return nil, fmt.Errorf("failed to find ticket availability: %w", err)
	}

	// Use the aggregate to enforce booking business rules
	if err := ticketAvailability.ReserveTickets(req.TicketsBooked); err != nil {
		s.logger.Warn().
			Err(err).
			Str("event_id", req.EventID.String()).
			Int("requested", req.TicketsBooked).
			Int("available", ticketAvailability.AvailableTickets).
			Msg("insufficient tickets")
		return nil, err
	}

	// Update the aggregate
	if err := s.ticketAvailabilityRepo.UpdateWithExecutor(ctx, tx, ticketAvailability); err != nil {
		s.logger.Error().
			Err(err).
			Str("event_id", req.EventID.String()).
			Msg("failed to update ticket availability")
		return nil, fmt.Errorf("failed to update ticket availability: %w", err)
	}

	booking, err := domain.NewBooking(req.EventID, req.UserID, req.TicketsBooked)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create booking domain object")
		return nil, fmt.Errorf("invalid booking data: %w", err)
	}

	if err := s.bookingRepo.CreateWithExecutor(ctx, tx, booking); err != nil {
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
