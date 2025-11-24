package domain

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID            uuid.UUID
	EventID       uuid.UUID
	UserID        uuid.UUID
	TicketsBooked int
	BookedAt      time.Time
}

func NewBooking(eventID, userID uuid.UUID, ticketsBooked int) (*Booking, error) {
	if ticketsBooked <= 0 {
		return nil, ErrInvalidTicketCount
	}

	return &Booking{
		ID:            uuid.New(),
		EventID:       eventID,
		UserID:        userID,
		TicketsBooked: ticketsBooked,
		BookedAt:      time.Now(),
	}, nil
}
