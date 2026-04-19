package http_response

import (
	"errors"
	"net/http"

	"github.com/egorkto/Chat-go/internal/domain"
)

func ErrorToHTTPCode(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalidArgument):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
