package users_storage

import (
	db_gorm_postgres "github.com/egorkto/Chat-go/internal/db/gorm/postgres"
)

type UsersStorage struct {
	db *db_gorm_postgres.DB
}

func New(db *db_gorm_postgres.DB) *UsersStorage {
	return &UsersStorage{
		db: db,
	}
}
