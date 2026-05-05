package users_storage

import (
	"context"
	"fmt"

	"github.com/egorkto/Chat-go/internal/domain"
	storage_postgres_gorm "github.com/egorkto/Chat-go/internal/storage/postgres/gorm"
)

func (s *UsersStorage) CreateUser(
	ctx context.Context,
	user domain.User,
	password string,
) (domain.User, error) {
	cancel := s.db.WithTimeoutContext(ctx)
	defer cancel()

	model := storage_postgres_gorm.UserModel{
		FullName: user.FullName(),
		Login:    user.Login(),
		Password: password,
	}

	err := s.db.Create(&model)
	if err != nil {
		return domain.User{}, fmt.Errorf("creating new user: %w", err)
	}
	domainUser := model.ToDomain()

	return domainUser, nil
}
