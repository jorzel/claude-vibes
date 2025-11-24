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

type EventHandler struct {
	service *app.EventService
	logger  zerolog.Logger
}

func NewEventHandler(service *app.EventService, logger zerolog.Logger) *EventHandler {
	return &EventHandler{
		service: service,
		logger:  logger.With().Str("handler", "event").Logger(),
	}
}

type CreateEventRequest struct {
	Name     string    `json:"name" validate:"required"`
	Date     time.Time `json:"date" validate:"required"`
	Location string    `json:"location" validate:"required"`
	Tickets  int       `json:"tickets" validate:"required,min=0"`
}

type EventResponse struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Date             time.Time `json:"date"`
	Location         string    `json:"location"`
	AvailableTickets int       `json:"available_tickets"`
	Tickets          int       `json:"tickets"`
}

func (h *EventHandler) CreateEvent(c echo.Context) error {
	var req CreateEventRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error().Err(err).Msg("failed to bind request")
		infrastructure.EventsCreated.WithLabelValues("error").Inc()
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}

	event, err := h.service.CreateEvent(c.Request().Context(), app.CreateEventRequest{
		Name:     req.Name,
		Date:     req.Date,
		Location: req.Location,
		Tickets:  req.Tickets,
	})
	if err != nil {
		infrastructure.EventsCreated.WithLabelValues("error").Inc()
		return handleError(c, err)
	}

	infrastructure.EventsCreated.WithLabelValues("success").Inc()
	return c.JSON(http.StatusCreated, EventResponse{
		ID:               event.ID.String(),
		Name:             event.Name,
		Date:             event.Date,
		Location:         event.Location,
		AvailableTickets: event.AvailableTickets,
		Tickets:          event.Tickets,
	})
}

func (h *EventHandler) GetEvent(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid event id"})
	}

	event, err := h.service.GetEvent(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, EventResponse{
		ID:               event.ID.String(),
		Name:             event.Name,
		Date:             event.Date,
		Location:         event.Location,
		AvailableTickets: event.AvailableTickets,
		Tickets:          event.Tickets,
	})
}

func (h *EventHandler) ListEvents(c echo.Context) error {
	events, err := h.service.ListEvents(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}

	response := make([]EventResponse, 0, len(events))
	for _, event := range events {
		response = append(response, EventResponse{
			ID:               event.ID.String(),
			Name:             event.Name,
			Date:             event.Date,
			Location:         event.Location,
			AvailableTickets: event.AvailableTickets,
			Tickets:          event.Tickets,
		})
	}

	return c.JSON(http.StatusOK, response)
}
