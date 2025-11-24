# Observability Standards

## Structured Logging

### Go
```go
import "github.com/rs/zerolog"

logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

// Log with structured fields
logger.Info().Str("order_id", Order.id).Msg("order created")

// Log errors with context
logger.Err(err).Msg("payment failed")

// Add context to logger
logger = logger.With("service", "order-service", "version", "1.0.0")
```

### Python
```python
import logging
import structlog

# Configure structlog
structlog.configure(
    processors=[
        structlog.stdlib.add_log_level,
        structlog.stdlib.add_logger_name,
        structlog.processors.TimeStamper(fmt="iso"),
        structlog.processors.JSONRenderer()
    ],
    wrapper_class=structlog.stdlib.BoundLogger,
    logger_factory=structlog.stdlib.LoggerFactory(),
)

logger = structlog.get_logger()

# Log with structured fields
logger.info(
    "order created",
    order_id=order.id,
    user_id=order.user_id,
    total=order.total,
    duration_ms=duration_ms,
)

# Log errors with context
logger.error(
    "payment failed",
    exc_info=True,
    order_id=order.id,
    amount=amount,
)

# Bind context
logger = logger.bind(service="order-service", version="1.0.0")
```

---

## Log Levels
- **ERROR**: System failures requiring immediate attention
- **WARN**: Degraded functionality, fallback engaged
- **INFO**: Key business events to track flow (order placed, user registered)
- **DEBUG**: Detailed troubleshooting (not in production)

---

## What to Log

### ✅ DO Log
```go
// Business events
logger.Info("order created", "order_id", id, "total", total)

// External API calls
logger.Info("payment processed",
    "provider", "stripe",
    "duration_ms", duration,
    "success", true,
)

// Errors with full context
logger.Error("database query failed",
    "error", err,
    "query", "INSERT INTO orders",
    "user_id", userID,
)
```

### ❌ DON'T Log
- Sensitive data (passwords, tokens, credit cards, PII)
- Inside tight loops (aggregate instead)
- Redundant info already captured in metrics

---

## Metrics

### Go (Prometheus)
```go
import "github.com/prometheus/client_golang/prometheus"

var (
    ordersCreated = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "orders_created_total",
            Help: "Total number of orders created",
        },
        []string{"status"},
    )

    orderDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "order_creation_duration_seconds",
            Help: "Order creation duration",
            Buckets: prometheus.DefBuckets,
        },
        []string{"endpoint"},
    )
)

// Usage
ordersCreated.WithLabelValues("success").Inc()
orderDuration.WithLabelValues("/api/orders").Observe(duration.Seconds())
```

### Python (Prometheus)
```python
from prometheus_client import Counter, Histogram

orders_created = Counter(
    'orders_created_total',
    'Total number of orders created',
    ['status']
)

order_duration = Histogram(
    'order_creation_duration_seconds',
    'Order creation duration',
    ['endpoint']
)

# Usage
orders_created.labels(status='success').inc()
order_duration.labels(endpoint='/api/orders').observe(duration)
```

---

## Tracing

### Go (OpenTelemetry)
```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

func CreateOrder(ctx context.Context, req *CreateOrderRequest) error {
    tracer := otel.Tracer("order-service")
    ctx, span := tracer.Start(ctx, "CreateOrder")
    defer span.End()

    // Add attributes
    span.SetAttributes(
        attribute.String("user_id", req.UserID),
        attribute.Int("item_count", len(req.Items)),
    )

    // Pass context to child operations
    err := s.repo.Save(ctx, order)
    if err != nil {
        span.RecordError(err)
        return err
    }

    return nil
}
```

### Python (OpenTelemetry)
```python
from opentelemetry import trace

tracer = trace.get_tracer(__name__)

def create_order(request: CreateOrderRequest) -> Order:
    with tracer.start_as_current_span("create_order") as span:
        span.set_attribute("user_id", request.user_id)
        span.set_attribute("item_count", len(request.items))

        try:
            order = self.repository.save(order)
            return order
        except Exception as e:
            span.record_exception(e)
            raise
```

---

## Health Checks

### Go
```go
func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    health := map[string]string{
        "status": "healthy",
    }

    // Check database
    if err := s.db.PingContext(ctx); err != nil {
        health["status"] = "unhealthy"
        health["database"] = "unreachable"
        w.WriteHeader(http.StatusServiceUnavailable)
    }

    json.NewEncoder(w).Encode(health)
}
```

### Python (FastAPI)
```python
@app.get("/health")
async def health_check():
    health = {"status": "healthy"}

    # Check database
    try:
        await db.execute("SELECT 1")
    except Exception as e:
        health["status"] = "unhealthy"
        health["database"] = "unreachable"
        return JSONResponse(
            status_code=503,
            content=health
        )

    return health
```


## Dashboards
- Create dashboards in Grafana to visualize key metrics:
  - Request rates, error rates, latencies
  - Business metrics (orders created, revenue)
  - System health (CPU, memory, DB connections)
- Use python libraries like `grafanalib` to define dashboards as code for versioning and reproducibility.