package chat_storage_postgres

import (
	"context"
	"fmt"

	"github.com/egorkto/Chat-go/internal/domain"
	storage_postgres_gorm "github.com/egorkto/Chat-go/internal/storage/postgres/gorm"
)

func (s ChatStorage) GetMessages(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.Message, error) {
	cancel := s.db.WithTimeout(ctx)
	defer cancel()

	msgLimit := -1
	msgOffset := -1

	var sliceSize int
	if limit != nil {
		msgLimit = *limit
		sliceSize = *limit
	}

	if offset != nil {
		msgOffset = *offset
	}

	models := make([]storage_postgres_gorm.MessageModel, sliceSize)

	err := s.db.Preload("User").
		Order("sent_at DESC").
		Limit(msgLimit).Offset(msgOffset).
		Find(&models).Error
	if err != nil {
		return []domain.Message{}, fmt.Errorf(
			"find with message models: %w",
			err,
		)
	}

	domainMsgs := make([]domain.Message, 0, len(models))

	for _, model := range models {
		domainMsgs = append(domainMsgs, model.ToDomain())
	}

	return domainMsgs, nil
}
