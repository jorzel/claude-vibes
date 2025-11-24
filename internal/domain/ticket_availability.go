package domain

import (
	"github.com/google/uuid"
)

// TicketAvailability is an aggregate that protects ticket reservation invariants
// It represents the consistency boundary for booking operations
type TicketAvailability struct {
	EventID          uuid.UUID
	AvailableTickets int
}

func NewTicketAvailability(eventID uuid.UUID, availableTickets int) (*TicketAvailability, error) {
	if availableTickets < 0 {
		return nil, ErrInvalidAvailableTickets
	}

	return &TicketAvailability{
		EventID:          eventID,
		AvailableTickets: availableTickets,
	}, nil
}

// ReserveTickets attempts to reserve the specified number of tickets
// This method enforces the invariant: AvailableTickets >= 0
func (ta *TicketAvailability) ReserveTickets(count int) error {
	if count <= 0 {
		return ErrInvalidTicketCount
	}

	if ta.AvailableTickets < count {
		return ErrInsufficientTickets
	}

	ta.AvailableTickets -= count
	return nil
}
