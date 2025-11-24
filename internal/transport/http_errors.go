package transport

import (
	"errors"
	"net/http"

	"github.com/jorzel/booking-service/internal/domain"
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func handleError(c echo.Context, err error) error {
	var notFoundErr *domain.NotFoundError
	var validationErr *domain.ValidationError
	var conflictErr *domain.ConflictError

	switch {
	case errors.As(err, &notFoundErr):
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case errors.As(err, &validationErr):
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	case errors.As(err, &conflictErr):
		return c.JSON(http.StatusConflict, ErrorResponse{Error: err.Error()})
	default:
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}
