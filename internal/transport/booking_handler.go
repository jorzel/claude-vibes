package transport

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jorzel/booking-service/internal/app"
	"github.com/jorzel/booking-service/internal/infrastructure"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type BookingHandler struct {
	service *app.BookingService
	logger  zerolog.Logger
}

func NewBookingHandler(service *app.BookingService, logger zerolog.Logger) *BookingHandler {
	return &BookingHandler{
		service: service,
		logger:  logger.With().Str("handler", "booking").Logger(),
	}
}

type CreateBookingRequest struct {
	EventID       string `json:"event_id" validate:"required"`
	UserID        string `json:"user_id" validate:"required"`
	TicketsBooked int    `json:"tickets_booked" validate:"required,min=1"`
}

type BookingResponse struct {
	ID            string    `json:"id"`
	EventID       string    `json:"event_id"`
	UserID        string    `json:"user_id"`
	TicketsBooked int       `json:"tickets_booked"`
	BookedAt      time.Time `json:"booked_at"`
}

func (h *BookingHandler) CreateBooking(c echo.Context) error {
	var req CreateBookingRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error().Err(err).Msg("failed to bind request")
		infrastructure.BookingsCreated.WithLabelValues("error").Inc()
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}

	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		infrastructure.BookingsCreated.WithLabelValues("error").Inc()
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid event_id"})
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		infrastructure.BookingsCreated.WithLabelValues("error").Inc()
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user_id"})
	}

	booking, err := h.service.CreateBooking(c.Request().Context(), app.CreateBookingRequest{
		EventID:       eventID,
		UserID:        userID,
		TicketsBooked: req.TicketsBooked,
	})
	if err != nil {
		infrastructure.BookingsCreated.WithLabelValues("error").Inc()
		return handleError(c, err)
	}

	infrastructure.BookingsCreated.WithLabelValues("success").Inc()
	infrastructure.TicketsBooked.Add(float64(booking.TicketsBooked))

	return c.JSON(http.StatusCreated, BookingResponse{
		ID:            booking.ID.String(),
		EventID:       booking.EventID.String(),
		UserID:        booking.UserID.String(),
		TicketsBooked: booking.TicketsBooked,
		BookedAt:      booking.BookedAt,
	})
}

func (h *BookingHandler) GetBooking(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid booking id"})
	}

	booking, err := h.service.GetBooking(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, BookingResponse{
		ID:            booking.ID.String(),
		EventID:       booking.EventID.String(),
		UserID:        booking.UserID.String(),
		TicketsBooked: booking.TicketsBooked,
		BookedAt:      booking.BookedAt,
	})
}
