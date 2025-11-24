package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTicketAvailability(t *testing.T) {
	tests := []struct {
		name             string
		eventID          uuid.UUID
		availableTickets int
		wantErr          bool
		errType          error
	}{
		{
			name:             "creates ticket availability with valid data",
			eventID:          uuid.New(),
			availableTickets: 100,
			wantErr:          false,
		},
		{
			name:             "creates ticket availability with zero tickets",
			eventID:          uuid.New(),
			availableTickets: 0,
			wantErr:          false,
		},
		{
			name:             "returns error for negative tickets",
			eventID:          uuid.New(),
			availableTickets: -10,
			wantErr:          true,
			errType:          ErrInvalidAvailableTickets,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			availability, err := NewTicketAvailability(tt.eventID, tt.availableTickets)

			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.errType))
				assert.Nil(t, availability)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, availability)
				assert.Equal(t, tt.eventID, availability.EventID)
				assert.Equal(t, tt.availableTickets, availability.AvailableTickets)
			}
		})
	}
}

func TestTicketAvailability_ReserveTickets(t *testing.T) {
	tests := []struct {
		name              string
		availableTickets  int
		requestedTickets  int
		wantErr           bool
		errType           error
		expectedRemaining int
	}{
		{
			name:              "reserves tickets successfully",
			availableTickets:  50,
			requestedTickets:  10,
			wantErr:           false,
			expectedRemaining: 40,
		},
		{
			name:              "reserves all available tickets",
			availableTickets:  20,
			requestedTickets:  20,
			wantErr:           false,
			expectedRemaining: 0,
		},
		{
			name:             "returns error when requesting more than available",
			availableTickets: 5,
			requestedTickets: 10,
			wantErr:          true,
			errType:          ErrInsufficientTickets,
		},
		{
			name:             "returns error for zero tickets",
			availableTickets: 100,
			requestedTickets: 0,
			wantErr:          true,
			errType:          ErrInvalidTicketCount,
		},
		{
			name:             "returns error for negative tickets",
			availableTickets: 100,
			requestedTickets: -5,
			wantErr:          true,
			errType:          ErrInvalidTicketCount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			availability := &TicketAvailability{
				EventID:          uuid.New(),
				AvailableTickets: tt.availableTickets,
			}

			err := availability.ReserveTickets(tt.requestedTickets)

			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.errType))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRemaining, availability.AvailableTickets)
			}
		})
	}
}
