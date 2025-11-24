1. Overview
The Booking Service is a REST API that manages event ticketing operations. Users can browse events and create bookings for available tickets.
Language: Go 1.24+
Framework: Echo (HTTP router)
Database: PostgreSQL

2. API Endpoints
- `POST /events`: Create a new event
- `GET /events`: List all events
- `GET /events/{id}`: Get event details
- `POST /bookings`: Create a new booking
- `GET /bookings/{id}`: Get booking details

3. Data Models
- Event
    - ID (UUID)
    - Name (string)
    - Date (timestamp)
    - Location (string)
    - AvailableTickets (int)
    - Tickets (int)

- Booking
    - ID (UUID)
    - EventID (UUID)
    - UserID (UUID)
    - TicketsBooked (int)
    - BookedAt (timestamp)

4. Database Schema
```sql
CREATE TABLE events (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    date TIMESTAMP NOT NULL,
    location VARCHAR(255) NOT NULL,
    available_tickets INT NOT NULL,
    tickets INT NOT NULL
);

CREATE TABLE bookings (
    id UUID PRIMARY KEY,
    event_id UUID REFERENCES events(id),
    user_id UUID NOT NULL,
    tickets_booked INT NOT NULL,
    booked_at TIMESTAMP NOT NULL
);
```

5. Business Logic
- When creating a booking, ensure that the number of tickets requested does not exceed the available tickets for the event.
- Upon successful booking, decrement the available tickets for the event accordingly.

6. Error Handling
- Return `400 Bad Request` for invalid input data.
- Return `409 Conflict` if trying to book more tickets than available.
- Return `404 Not Found` if the specified event or booking does not exist.
- Return `500 Internal Server Error` for unexpected server errors.

7. Observability
- logging and metrics
- tracing can be added later

8. Testing
- Unit tests for business logic
- Integration tests for API endpoints
- Use testcontainers for e2e tests with a real PostgreSQL instance

