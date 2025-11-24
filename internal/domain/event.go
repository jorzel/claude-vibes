package domain

import (
	"time"

	"github.com/google/uuid"
)

// Event is a data container for event metadata
// It does not contain booking business logic - that is handled by TicketAvailability aggregate
type Event struct {
	ID       uuid.UUID
	Name     string
	Date     time.Time
	Location string
	Tickets  int // Total tickets (immutable reference)
}

func NewEvent(name, location string, date time.Time, tickets int) (*Event, error) {
	if tickets < 0 {
		return nil, ErrInvalidAvailableTickets
	}

	return &Event{
		ID:       uuid.New(),
		Name:     name,
		Date:     date,
		Location: location,
		Tickets:  tickets,
	}, nil
}
