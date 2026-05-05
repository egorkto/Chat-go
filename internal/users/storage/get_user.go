package users_storage

import (
	"context"
	"fmt"

	"github.com/egorkto/Chat-go/internal/domain"
	storage_postgres_gorm "github.com/egorkto/Chat-go/internal/storage/postgres/gorm"
)

func (s *UsersStorage) GetUserByLogin(
	ctx context.Context,
	login string,
) (domain.User, string, error) {
	cancel := s.db.WithTimeoutContext(ctx)
	defer cancel()

	var userModel storage_postgres_gorm.UserModel

	err := s.db.First(&userModel, "login = ?", login)
	if err != nil {
		return domain.User{}, "", fmt.Errorf(
			"recieving user by login and password: %w",
			err,
		)
	}

	user := userModel.ToDomain()

	return user, userModel.Password, nil
}

func (s *UsersStorage) GetUserByID(
	ctx context.Context,
	id int,
) (domain.User, error) {
	cancel := s.db.WithTimeoutContext(ctx)
	defer cancel()

	var userModel storage_postgres_gorm.UserModel

	err := s.db.First(&userModel, "id = ?", id)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"recieving user by id: %w",
			err,
		)
	}

	user := userModel.ToDomain()

	return user, nil
}
