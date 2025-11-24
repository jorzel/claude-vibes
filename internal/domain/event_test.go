package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEvent(t *testing.T) {
	tests := []struct {
		name     string
		evtName  string
		location string
		date     time.Time
		tickets  int
		wantErr  bool
		errType  error
	}{
		{
			name:     "creates event with valid data",
			evtName:  "Concert Night",
			location: "Madison Square Garden",
			date:     time.Now().Add(24 * time.Hour),
			tickets:  100,
			wantErr:  false,
		},
		{
			name:     "creates event with zero tickets",
			evtName:  "Free Event",
			location: "Central Park",
			date:     time.Now().Add(48 * time.Hour),
			tickets:  0,
			wantErr:  false,
		},
		{
			name:     "returns error for negative tickets",
			evtName:  "Invalid Event",
			location: "Somewhere",
			date:     time.Now(),
			tickets:  -10,
			wantErr:  true,
			errType:  ErrInvalidAvailableTickets,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := NewEvent(tt.evtName, tt.location, tt.date, tt.tickets)

			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.errType))
				assert.Nil(t, event)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, event)
				assert.NotEqual(t, "", event.ID.String())
				assert.Equal(t, tt.evtName, event.Name)
				assert.Equal(t, tt.location, event.Location)
				assert.Equal(t, tt.tickets, event.Tickets)
			}
		})
	}
}
