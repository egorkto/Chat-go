package storage_postgres_gorm

import (
	"context"
	"fmt"

	storage_postgres "github.com/egorkto/Chat-go/internal/storage/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
	cfg storage_postgres.Config
}

func New(cfg storage_postgres.Config) (*DB, error) {
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

	return &DB{
		DB:  db,
		cfg: cfg,
	}, nil
}

func (db *DB) WithTimeoutContext(ctx context.Context) context.CancelFunc {
	timeoutCtx, cancel := context.WithTimeout(ctx, db.cfg.Timeout)

	db.DB = db.DB.WithContext(timeoutCtx)

	return cancel
}

func (db *DB) Create(dest interface{}) error {
	result := db.DB.Create(dest)
	return MapError(result.Error)
}

func (db *DB) Where(query interface{}, args ...interface{}) error {
	result := db.DB.Where(query, args...)
	return MapError(result.Error)
}

func (db *DB) First(dest interface{}, query interface{}, args ...interface{}) error {
	result := db.DB.Where(query, args...).First(dest)
	return MapError(result.Error)
}
