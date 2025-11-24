package domain

import "fmt"

var (
	ErrEventNotFound          = &NotFoundError{Entity: "event"}
	ErrBookingNotFound        = &NotFoundError{Entity: "booking"}
	ErrInsufficientTickets    = &ConflictError{Message: "insufficient tickets available"}
	ErrInvalidTicketCount     = &ValidationError{Field: "tickets_booked", Message: "must be greater than 0"}
	ErrInvalidAvailableTickets = &ValidationError{Field: "available_tickets", Message: "cannot be negative"}
)

type NotFoundError struct {
	Entity string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Entity)
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict: %s", e.Message)
}
