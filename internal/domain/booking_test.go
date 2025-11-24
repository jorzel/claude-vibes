package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBooking(t *testing.T) {
	eventID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name          string
		eventID       uuid.UUID
		userID        uuid.UUID
		ticketsBooked int
		wantErr       bool
		errType       error
	}{
		{
			name:          "creates booking with valid data",
			eventID:       eventID,
			userID:        userID,
			ticketsBooked: 3,
			wantErr:       false,
		},
		{
			name:          "creates booking with single ticket",
			eventID:       eventID,
			userID:        userID,
			ticketsBooked: 1,
			wantErr:       false,
		},
		{
			name:          "returns error for zero tickets",
			eventID:       eventID,
			userID:        userID,
			ticketsBooked: 0,
			wantErr:       true,
			errType:       ErrInvalidTicketCount,
		},
		{
			name:          "returns error for negative tickets",
			eventID:       eventID,
			userID:        userID,
			ticketsBooked: -2,
			wantErr:       true,
			errType:       ErrInvalidTicketCount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			booking, err := NewBooking(tt.eventID, tt.userID, tt.ticketsBooked)

			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.errType))
				assert.Nil(t, booking)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, booking)
				assert.NotEqual(t, "", booking.ID.String())
				assert.Equal(t, tt.eventID, booking.EventID)
				assert.Equal(t, tt.userID, booking.UserID)
				assert.Equal(t, tt.ticketsBooked, booking.TicketsBooked)
				assert.False(t, booking.BookedAt.IsZero())
			}
		})
	}
}
