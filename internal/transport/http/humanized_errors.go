package transport_http

import (
	"errors"

	"github.com/egorkto/Chat-go/internal/domain"
)

func GetHumanized(err error) string {
	for k, v := range humanizedErrors {
		if errors.Is(err, k) {
			return v
		}
	}

	return "Internal server error"
}

var invalidCredsErrorText = "invalid login or password"

var humanizedErrors = map[error]string{
	domain.ErrPasswordNotMatch:    invalidCredsErrorText,
	domain.ErrUserNotFoundByLogin: invalidCredsErrorText,
	domain.ErrUserNotFoundById:    "user not found",
	domain.ErrLoginDuplicated:     "such login already exists",
	domain.ErrUnauthorized:        "unauthorized request",
	domain.ErrValidationFailed:    "validation error",
}
