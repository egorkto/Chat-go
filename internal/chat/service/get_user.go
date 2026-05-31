package chat_service

import (
	"context"
	"fmt"

	"github.com/egorkto/Chat-go/internal/domain"
)

func (s ChatService) GetUser(ctx context.Context, id int) (domain.User, error) {
	user, err := s.usersStorage.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}
