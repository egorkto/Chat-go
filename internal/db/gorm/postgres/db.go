package db_gorm_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/egorkto/Chat-go/internal/db"
	"github.com/egorkto/Chat-go/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
	cfg db.Config
}

func New(cfg db.Config) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	return &DB{
		DB:  db,
		cfg: cfg,
	}, nil
}

func (db *DB) WithTimeoutContextBasedOn(ctx context.Context) context.CancelFunc {
	timeoutCtx, cancel := context.WithTimeout(ctx, db.cfg.Timeout)

	db.DB = db.DB.WithContext(timeoutCtx)

	return cancel
}

func (db *DB) Create(dest interface{}) error {
	result := db.DB.Create(dest)
	err := result.Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrCheckConstraintViolated):
			return domain.ErrConflict
		case errors.Is(err, gorm.ErrDuplicatedKey):
			return domain.ErrInvalidArgument
		case errors.Is(err, gorm.ErrRecordNotFound):
			return domain.ErrNotFound
		default:
			return err
		}
	}
	return nil
}
