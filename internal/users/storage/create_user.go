package users_storage

import (
	"context"
	"fmt"

	db_gorm "github.com/egorkto/Chat-go/internal/db/gorm"
	"github.com/egorkto/Chat-go/internal/domain"
)

func (s *UsersStorage) CreateUser(
	ctx context.Context,
	user domain.User,
	password string,
) (domain.User, error) {
	cancel := s.db.WithTimeoutContextBasedOn(ctx)
	defer cancel()

	model := db_gorm.UserModel{
		FullName: user.FullName(),
		Password: password,
	}

	err := s.db.Create(&model)
	if err != nil {
		return domain.User{}, fmt.Errorf("creating new user: %w", err)
	}
	domainUser := model.ToDomain()

	return domainUser, nil
}
