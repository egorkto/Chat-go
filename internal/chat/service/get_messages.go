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
			domain.NewValidationError(map[string]string{
				"limit": "limit must be non-negative",
			}),
		)
	}

	if offset != nil && *offset < 0 {
		return []domain.Message{}, fmt.Errorf(
			"offset must be non-negative: %w",
			domain.NewValidationError(map[string]string{
				"offset": "offset must be non-negative",
			}),
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
