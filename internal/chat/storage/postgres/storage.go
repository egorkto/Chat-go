package chat_storage_postgres

import (
	storage_postgres_gorm "github.com/egorkto/Chat-go/internal/storage/postgres/gorm"
)

type ChatStorage struct {
	db *storage_postgres_gorm.GormDB
}

func New(db *storage_postgres_gorm.GormDB) *ChatStorage {
	return &ChatStorage{
		db: db,
	}
}
