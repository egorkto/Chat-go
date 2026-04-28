package db_gorm

import "context"

type DB interface {
	WithTimeoutContext(ctx context.Context) context.CancelFunc
	Create(dest interface{}) error
}
