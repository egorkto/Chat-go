package storage_postgres_gorm

import (
	"context"
	"fmt"
	"time"

	storage_postgres "github.com/egorkto/Chat-go/internal/storage/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	*gorm.DB
	timeout time.Duration
}

func New(cfg storage_postgres.Config) (*GormDB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	return &GormDB{
		DB:      db,
		timeout: cfg.Timeout,
	}, nil
}

func (db *GormDB) SetTimeout(timeout time.Duration) error {
	db.timeout = timeout
	return nil
}

func (db GormDB) WithTimeout(ctx context.Context) context.CancelFunc {
	timeoutCtx, cancel := context.WithTimeout(ctx, db.timeout)
	db.DB = db.DB.WithContext(timeoutCtx)
	return cancel
}
