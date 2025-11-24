package infrastructure

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	EventsCreated = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "booking_service_events_created_total",
			Help: "Total number of events created",
		},
		[]string{"status"},
	)

	BookingsCreated = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "booking_service_bookings_created_total",
			Help: "Total number of bookings created",
		},
		[]string{"status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "booking_service_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	TicketsBooked = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "booking_service_tickets_booked_total",
			Help: "Total number of tickets booked",
		},
	)

	PostgresQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "booking_service_postgres_queries_total",
			Help: "Total number of Postgres queries executed",
		},
		[]string{"operation", "status"},
	)

	PostgresQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "booking_service_postgres_query_duration_seconds",
			Help:    "Postgres query duration in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
		},
		[]string{"operation"},
	)
)
