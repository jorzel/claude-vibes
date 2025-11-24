package domain

import (
	"context"

	"github.com/google/uuid"
)

type EventRepository interface {
	Create(ctx context.Context, event *Event) error
	FindByID(ctx context.Context, id uuid.UUID) (*Event, error)
	FindAll(ctx context.Context) ([]*Event, error)
	Update(ctx context.Context, event *Event) error
}

type BookingRepository interface {
	Create(ctx context.Context, booking *Booking) error
	FindByID(ctx context.Context, id uuid.UUID) (*Booking, error)
}
