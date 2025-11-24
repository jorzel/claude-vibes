# Claude Vibes

Vibe engineering with Claude - demonstrating clean architecture and test-driven development.

## Projects

### Booking Service

A REST API for managing event ticketing operations built with Go, following Domain-Driven Design principles and clean architecture patterns.

#### Features

- **Event Management**: Create and browse events with ticket availability
- **Booking System**: Book tickets with automatic inventory management
- **Concurrency Safety**: Handles concurrent bookings with database transactions
- **Observability**: Structured logging with zerolog and Prometheus metrics
- **Testing**: Comprehensive unit and integration tests with testcontainers

#### Architecture

The project follows a clean architecture pattern with clear separation of concerns:

```
cmd/server/          # Application entry point
internal/
  ├── domain/        # Business logic and entities
  ├── app/           # Application services
  ├── infrastructure/# External dependencies (database, metrics)
  └── transport/     # HTTP handlers and routing
tests/               # Integration tests
```

#### Tech Stack

- **Language**: Go 1.23+
- **Web Framework**: Echo
- **Database**: PostgreSQL 16
- **Logging**: zerolog
- **Metrics**: Prometheus
- **Testing**: testify, testcontainers

#### API Endpoints

**Events**
- `POST /events` - Create a new event
- `GET /events` - List all events
- `GET /events/{id}` - Get event details

**Bookings**
- `POST /bookings` - Create a new booking
- `GET /bookings/{id}` - Get booking details

**Health & Metrics**
- `GET /health` - Health check endpoint
- `GET /metrics` - Prometheus metrics

#### Getting Started

**Prerequisites**
- Go 1.23+
- Docker and Docker Compose
- Make (optional)

**Option 1: Using Docker Compose (Recommended)**

Start all services (PostgreSQL, booking service, Prometheus, Grafana):
```bash
docker-compose up -d
```

This will:
- Start PostgreSQL with automatic migration execution
- Start the booking service on `http://localhost:8080`
- Start Prometheus on `http://localhost:9090`
- Start Grafana on `http://localhost:3000` (admin/admin)

**Option 2: Local Development**

For running the server locally with containerized PostgreSQL:

1. Install dependencies:
```bash
make deps
```

2. Start PostgreSQL only:
```bash
make docker-up
```

3. Run database migrations (requires `psql` installed locally):
```bash
make migrate
```

4. Start the server:
```bash
make run
```

The server will start on `http://localhost:8080`.

**Testing**

```bash
make test                 # Run all tests
make unit-test            # Run unit tests
make integration-test     # Run integration tests
make coverage             # Generate coverage report
make bench                # Run benchmarks
```

**Additional Commands**

```bash
make help                 # Show all available commands
make build                # Build the binary
make lint                 # Run linter
make fmt                  # Format code
make dashboard            # Generate Grafana dashboard
make dashboard-deps       # Install dashboard dependencies
make openapi-validate     # Validate OpenAPI specification
```

#### Usage Examples

**Create an Event**
```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Rock Concert",
    "date": "2025-12-31T20:00:00Z",
    "location": "Madison Square Garden",
    "tickets": 1000
  }'
```

**Create a Booking**
```bash
curl -X POST http://localhost:8080/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": "EVENT_UUID",
    "user_id": "USER_UUID",
    "tickets_booked": 2
  }'
```

#### API Documentation

The API is fully documented using OpenAPI 3.0 specification:

- **OpenAPI Spec**: `openapi.yaml` - Complete API documentation with request/response schemas
- **Validation**: Run `make openapi-validate` to validate the specification

You can use the OpenAPI spec with tools like:
- Swagger UI for interactive API documentation
- Postman for API testing
- Code generation tools (openapi-generator)

#### Monitoring & Dashboards

The service includes comprehensive Grafana dashboards for monitoring:

**Available Dashboards:**
- `dashboards/booking_service.json` - Pre-built dashboard (ready to import)
- `dashboards/booking_service_dashboard.py` - Dashboard as code using grafanalib

**Dashboard Features:**
- Business metrics (events, bookings, tickets)
- HTTP performance (request rate, latency percentiles)
- Error tracking (status codes, error rates)
- System metrics (goroutines, memory)

**Setup:**

1. Install dashboard dependencies (optional, for Python generation):
```bash
make dashboard-deps
```

2. Generate dashboard from Python code (optional):
```bash
make dashboard
```

3. Import into Grafana:
   - Open Grafana → Dashboards → Import
   - Upload `dashboards/booking_service.json`
   - Select Prometheus datasource

See `dashboards/README.md` for detailed documentation.

#### Configuration

Environment variables:
- `DB_HOST` - Database host (default: localhost)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: postgres)
- `DB_NAME` - Database name (default: booking_service)
- `DB_SSLMODE` - SSL mode (default: disable)
- `PORT` - Server port (default: 8080)

## Development Guidelines

See `.claude/` directory for detailed guidelines:
- `CLAUDE_CONTEXT.md` - Project overview and principles
- `STYLE_GUIDE.md` - Go coding conventions
- `TESTING_SPEC.md` - Testing standards
- `OBSERVABILITY.md` - Logging and metrics standards

## License

See LICENSE file for details.
