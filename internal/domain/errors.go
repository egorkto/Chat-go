package domain

import (
	"errors"
)

var (
	ErrValidationFailed    = errors.New("validation failed")
	ErrLoginDuplicated     = errors.New("login duplicated")
	ErrUserNotFoundByLogin = errors.New("user not found by login")
	ErrUserNotFoundById    = errors.New("user not found by id")
	ErrPasswordNotMatch    = errors.New("password does not match")
	ErrUnauthorized        = errors.New("unauthorized")
)

type ValidationError struct {
	Errors map[string]string
}

func NewValidationError(errs map[string]string) ValidationError {
	return ValidationError{
		Errors: errs,
	}
}

func (e ValidationError) Error() string {
	return ErrValidationFailed.Error()
}

func (e ValidationError) Is(target error) bool {
	return target == ErrValidationFailed
}
