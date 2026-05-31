package transport_http

import (
	"errors"
	"net/http"

	"github.com/egorkto/Chat-go/internal/domain"
)

func ErrorToHTTPCode(err error) int {
	switch {
	case errors.Is(err, domain.ErrValidationFailed),
		errors.Is(err, domain.ErrUserNotFoundByLogin),
		errors.Is(err, domain.ErrPasswordNotMatch):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrUserNotFoundById):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrLoginDuplicated):
		return http.StatusConflict
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
