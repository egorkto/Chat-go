package storage_postgres_gorm

import (
	"fmt"
	"strings"

	"github.com/egorkto/Chat-go/internal/domain"
)

func MapConstrainedError(err error) error {
	if err == nil {
		return nil
	}

	var mappedErr error
	if strings.Contains(err.Error(), "chk_full_name_length") {
		mappedErr = domain.NewValidationError(map[string]string{
			"full_name": "invalid length",
		})
	} else if strings.Contains(err.Error(), "chk_login_length") {
		mappedErr = domain.NewValidationError(map[string]string{
			"login": "invalid length",
		})
	} else if strings.Contains(err.Error(), "chk_password_length") {
		mappedErr = domain.NewValidationError(map[string]string{
			"password": "invalid length",
		})
	} else if strings.Contains(err.Error(), "uq_login") {
		mappedErr = domain.ErrLoginDuplicated
	}

	if mappedErr != nil {
		return fmt.Errorf("%s:%w", err.Error(), mappedErr)
	}

	return err
}
