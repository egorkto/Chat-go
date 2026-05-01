package db_gorm_postgres

import (
	"errors"
	"fmt"

	db_gorm "github.com/egorkto/Chat-go/internal/db/gorm"
	"gorm.io/gorm"
)

func MapError(err error) error {
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrCheckConstraintViolated):
			return fmt.Errorf("%s:%w", err.Error(), db_gorm.ErrCheckConstraintViolated)
		case errors.Is(err, gorm.ErrDuplicatedKey):
			return fmt.Errorf("%s:%w", err.Error(), db_gorm.ErrDuplicatedKey)
		case errors.Is(err, gorm.ErrRecordNotFound):
			return fmt.Errorf("%s:%w", err, db_gorm.ErrRecordNotFound)
		default:
			return err
		}
	}
	return nil
}
