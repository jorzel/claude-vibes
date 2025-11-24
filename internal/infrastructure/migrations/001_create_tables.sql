-- Create events table
CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    date TIMESTAMP NOT NULL,
    location VARCHAR(255) NOT NULL,
    available_tickets INT NOT NULL,
    tickets INT NOT NULL,
    CONSTRAINT available_tickets_non_negative CHECK (available_tickets >= 0),
    CONSTRAINT tickets_non_negative CHECK (tickets >= 0)
);

-- Create bookings table
CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY,
    event_id UUID NOT NULL REFERENCES events(id),
    user_id UUID NOT NULL,
    tickets_booked INT NOT NULL,
    booked_at TIMESTAMP NOT NULL,
    CONSTRAINT tickets_booked_positive CHECK (tickets_booked > 0)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_bookings_event_id ON bookings(event_id);
CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_events_date ON events(date);
