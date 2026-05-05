package users_service

import (
	"context"
	"fmt"

	"github.com/egorkto/Chat-go/internal/domain"
)

func (s *UsersService) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	user, err := s.storage.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("getting user: %w", err)
	}

	return user, nil
}
