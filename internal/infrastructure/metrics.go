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
)
