package constant

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrValidation maps to http.StatusBadRequest
type ErrValidation struct {
	Violations error
}

func (e *ErrValidation) Error() string {
	return fmt.Sprintf("there are validation errors: %v", e.Violations)
}

// internal app specific errors to be used
var (
	ErrConflict              = errors.New("the entity is in a conflict")
	ErrDatabase              = errors.New("encountered a database error")
	ErrDBNoSuchEntity        = errors.New("no such entity")
	ErrDBEntityAlreadyExists = errors.New("entity with same key already exists")
	ErrUnauthorized          = errors.New("user is unauthorized")
	ErrNoPermission          = errors.New("user has insufficient permissions")
)

// ErrToHTTPCode maps the application errors to specific http status codes.
// By default, it maps to http.StatusInternalServerError
func ErrToHTTPCode(err error) int {
	var errValidation *ErrValidation
	switch {
	case errors.As(err, &errValidation):
		return http.StatusBadRequest
	case errors.Is(err, ErrConflict):
		return http.StatusConflict
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrNoPermission):
		return http.StatusForbidden
	case errors.Is(err, ErrDatabase):
		return http.StatusInternalServerError
	case errors.Is(err, ErrDBNoSuchEntity):
		return http.StatusNotFound
	case errors.Is(err, ErrDBEntityAlreadyExists):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
