package chat_storage_postgres

import (
	"context"
	"fmt"

	"github.com/egorkto/Chat-go/internal/domain"
	storage_postgres_gorm "github.com/egorkto/Chat-go/internal/storage/postgres/gorm"
)

func (s ChatStorage) CreateMessage(
	ctx context.Context,
	msg domain.Message,
) (domain.Message, error) {
	cancel := s.db.WithTimeout(ctx)
	defer cancel()

	model := storage_postgres_gorm.MessageModel{
		Text:   msg.Text,
		UserID: msg.Sender.ID(),
		SentAt: msg.SendTime,
	}

	err := s.db.Create(&model).Error
	if err != nil {
		return domain.Message{}, fmt.Errorf(
			"create with message model: %w",
			err,
		)
	}

	err = s.db.Preload("User").First(&model, model.ID).Error
	if err != nil {
		return domain.Message{}, fmt.Errorf(
			"preload user: %w",
			err,
		)
	}

	domainMsg := model.ToDomain()

	return domainMsg, nil
}
