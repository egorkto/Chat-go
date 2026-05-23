package chat_service

import (
	"context"
	"fmt"

	"github.com/egorkto/Chat-go/internal/domain"
)

func (s ChatService) GetMessages(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.Message, error) {
	if limit != nil && *limit < 0 {
		return []domain.Message{}, fmt.Errorf(
			"limit must be non-negative: %w",
			domain.ErrInvalidArgument,
		)
	}

	if offset != nil && *offset < 0 {
		return []domain.Message{}, fmt.Errorf(
			"offset must be non-negative: %w",
			domain.ErrInvalidArgument,
		)
	}

	msgs, err := s.msgStorage.GetMessages(ctx, limit, offset)
	if err != nil {
		return []domain.Message{}, fmt.Errorf(
			"get messages: %w",
			err,
		)
	}

	return msgs, nil
}
