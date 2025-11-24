package domain

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID               uuid.UUID
	Name             string
	Date             time.Time
	Location         string
	AvailableTickets int
	Tickets          int
}

func NewEvent(name, location string, date time.Time, tickets int) (*Event, error) {
	if tickets < 0 {
		return nil, ErrInvalidAvailableTickets
	}

	return &Event{
		ID:               uuid.New(),
		Name:             name,
		Date:             date,
		Location:         location,
		AvailableTickets: tickets,
		Tickets:          tickets,
	}, nil
}

func (e *Event) ReserveTickets(count int) error {
	if count <= 0 {
		return ErrInvalidTicketCount
	}

	if e.AvailableTickets < count {
		return ErrInsufficientTickets
	}

	e.AvailableTickets -= count
	return nil
}
