package transport

import (
	"database/sql"
	"net/http"

	"github.com/jorzel/booking-service/internal/app"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

func NewRouter(
	eventService *app.EventService,
	bookingService *app.BookingService,
	db *sql.DB,
	logger zerolog.Logger,
) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.RequestID())
	e.Use(LoggingMiddleware(logger))
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
