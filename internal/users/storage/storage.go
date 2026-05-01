package users_storage

import (
	db_gorm "github.com/egorkto/Chat-go/internal/db/gorm"
)

type UsersStorage struct {
	db db_gorm.DB
}

func New(db db_gorm.DB) *UsersStorage {
	return &UsersStorage{
		db: db,
	}
}
