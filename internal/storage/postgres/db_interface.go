package storage_postgres

import "context"

type DB interface {
	WithTimeoutContext(ctx context.Context) context.CancelFunc
	Create(dest interface{}) error
	First(dest interface{}, query interface{}, args ...interface{}) error
}
