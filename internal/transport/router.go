package transport

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jorzel/booking-service/internal/app"
	"github.com/jorzel/booking-service/internal/infrastructure"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

func NewRouter(
	eventService *app.EventService,
	bookingService *app.BookingService,
	db infrastructure.DBClient,
	logger zerolog.Logger,
) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.RequestID())
	e.Use(LoggingMiddleware(logger))
	e.Use(MetricsMiddleware())
	e.Use(middleware.Recover())

	eventHandler := NewEventHandler(eventService, logger)
	bookingHandler := NewBookingHandler(bookingService, logger)

	e.POST("/events", eventHandler.CreateEvent)
	e.GET("/events", eventHandler.ListEvents)
	e.GET("/events/:id", eventHandler.GetEvent)

	e.POST("/bookings", bookingHandler.CreateBooking)
	e.GET("/bookings/:id", bookingHandler.GetBooking)

	e.GET("/health", func(c echo.Context) error {
		if err := db.PingContext(c.Request().Context()); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"status":   "unhealthy",
				"database": "unreachable",
			})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	return e
}

func LoggingMiddleware(logger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			logger.Info().
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Str("request_id", req.Header.Get(echo.HeaderXRequestID)).
				Msg("incoming request")

			err := next(c)

			logger.Info().
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Int("status", res.Status).
				Str("request_id", req.Header.Get(echo.HeaderXRequestID)).
				Msg("request completed")

			return err
		}
	}
}

func MetricsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip metrics for /metrics endpoint to avoid self-instrumentation
			if c.Request().URL.Path == "/metrics" {
				return next(c)
			}

			start := time.Now()
			err := next(c)
			duration := time.Since(start).Seconds()

			status := c.Response().Status
			method := c.Request().Method
			path := c.Path()

			// Record HTTP request duration
			infrastructure.HTTPRequestDuration.WithLabelValues(
				method,
				path,
				strconv.Itoa(status),
			).Observe(duration)

			return err
		}
	}
}
