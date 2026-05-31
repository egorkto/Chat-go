package users_storage_postgres

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
	cancel := s.db.WithTimeout(ctx)
	defer cancel()

	model := storage_postgres_gorm.UserModel{
		FullName: user.FullName(),
		Login:    user.Login(),
		Password: password,
	}

	err := s.db.DB.Create(&model).Error
	if err != nil {
		mapped := storage_postgres_gorm.MapConstrainedError(err)
		return domain.User{}, fmt.Errorf(
			"create user: %w",
			mapped,
		)
	}
	domainUser := model.ToDomain()

	return domainUser, nil
}
