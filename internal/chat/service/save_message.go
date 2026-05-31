package chat_service

import (
	"context"
	"fmt"

	"github.com/egorkto/Chat-go/internal/domain"
)

func (s ChatService) SaveMessage(
	ctx context.Context,
	msg domain.Message,
) error {
	msg, err := s.msgStorage.CreateMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("create message: %w", err)
	}

	return nil
}
