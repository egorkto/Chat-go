package transport_http

import (
	"errors"

	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/go-playground/validator/v10"
)

func ParseValidateError(err error) error {
	if err == nil {
		return nil
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errs := make(map[string]string)

		for _, fe := range ve {
			fieldName := fe.Field()
			errs[fieldName] = fe.Error()
		}

		return domain.NewValidationError(errs)
	}

	return err
}
