package users_storage_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/egorkto/Chat-go/internal/domain"
	storage_postgres_gorm "github.com/egorkto/Chat-go/internal/storage/postgres/gorm"
	"gorm.io/gorm"
)

func (s *UsersStorage) GetUserByLogin(
	ctx context.Context,
	login string,
) (domain.User, string, error) {
	cancel := s.db.WithTimeout(ctx)
	defer cancel()

	var userModel storage_postgres_gorm.UserModel

	err := s.db.First(&userModel, "login = ?", login).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, "", fmt.Errorf(
				"user record not found, %s: %w",
				err.Error(),
				domain.ErrNotFound,
			)
		}
		return domain.User{}, "", fmt.Errorf(
			"get user by login: %w",
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
	cancel := s.db.WithTimeout(ctx)
	defer cancel()

	var userModel storage_postgres_gorm.UserModel

	err := s.db.First(&userModel, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, fmt.Errorf(
				"user record not found, %s: %w",
				err.Error(),
				domain.ErrNotFound,
			)
		} else {
			return domain.User{}, fmt.Errorf(
				"recieving user by id: %w",
				err,
			)
		}
	}

	user := userModel.ToDomain()

	return user, nil
}
