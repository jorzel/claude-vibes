-- Create ticket_availability table as a separate aggregate
-- This protects booking business rules and eliminates lock contention with event management
CREATE TABLE IF NOT EXISTS ticket_availability (
    event_id UUID PRIMARY KEY REFERENCES events(id),
    available_tickets INT NOT NULL,
    CONSTRAINT available_tickets_non_negative CHECK (available_tickets >= 0)
);

-- Migrate existing available_tickets data from events to ticket_availability
INSERT INTO ticket_availability (event_id, available_tickets)
SELECT id, available_tickets FROM events
ON CONFLICT (event_id) DO NOTHING;

-- Remove available_tickets from events table as it's now in the aggregate
ALTER TABLE events DROP COLUMN IF EXISTS available_tickets;

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_ticket_availability_event_id ON ticket_availability(event_id);
