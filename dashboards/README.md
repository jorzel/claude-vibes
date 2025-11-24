# Grafana Dashboards

This directory contains Grafana dashboard definitions for monitoring the Booking Service.

## Available Dashboards

### Booking Service Monitoring Dashboard

A comprehensive dashboard that tracks:

**Business Metrics:**
- Total events created
- Total bookings created
- Total tickets booked
- Events and bookings creation rate over time

**Performance Metrics:**
- HTTP request rate by endpoint
- Response time percentiles (p50, p95, p99)
- Request latency distribution

**Error Tracking:**
- HTTP error rate by status code (4xx, 5xx)
- Error rate percentage
- Failed event/booking operations

**System Metrics:**
- Go goroutines
- Memory usage (allocated, heap)

## Using the Dashboards

### Option 1: Import Pre-generated JSON

The `booking_service.json` file contains a ready-to-use dashboard that can be imported directly into Grafana:

1. Open Grafana web interface
2. Navigate to Dashboards → Import
3. Upload `dashboards/booking_service.json`
4. Select your Prometheus datasource
5. Click Import

### Option 2: Generate from Python Code

For customization, you can generate the dashboard from the Python definition:

1. Install dependencies:
```bash
cd dashboards
pip install -r requirements.txt
```

2. Generate the dashboard JSON:
```bash
python booking_service_dashboard.py > booking_service_generated.json
```

3. Import the generated JSON into Grafana

### Option 3: Use Makefile

From the project root:

```bash
# Generate dashboard JSON from Python
make dashboard

# Install dashboard dependencies
make dashboard-deps
```

## Customization

To customize the dashboard:

1. Edit `booking_service_dashboard.py`
2. Modify panels, metrics, or add new visualizations
3. Regenerate the JSON file
4. Re-import into Grafana

## Metrics Reference

The dashboard uses these Prometheus metrics:

- `booking_service_events_created_total{status}` - Counter of events created
- `booking_service_bookings_created_total{status}` - Counter of bookings created
- `booking_service_tickets_booked_total` - Counter of tickets booked
- `booking_service_http_request_duration_seconds` - Histogram of request durations
- `go_goroutines` - Number of goroutines
- `go_memstats_alloc_bytes` - Bytes allocated
- `go_memstats_heap_alloc_bytes` - Heap bytes allocated

## Dashboard Preview

The dashboard is organized into rows:

1. **Overview** - High-level stats (single stat panels)
2. **Business Metrics** - Event and booking creation trends
3. **HTTP Performance** - Request rates and response times
4. **Error Rates** - Error tracking and percentages
5. **System Metrics** - Go runtime metrics

## Requirements

- Grafana 7.0+
- Prometheus datasource configured in Grafana
- Booking Service running and exposing metrics at `/metrics`
