package users_storage

import (
	storage_postgres "github.com/egorkto/Chat-go/internal/storage/postgres"
)

type UsersStorage struct {
	db storage_postgres.DB
}

func New(db storage_postgres.DB) *UsersStorage {
	return &UsersStorage{
		db: db,
	}
}
