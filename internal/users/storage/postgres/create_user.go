package users_storage_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/egorkto/Chat-go/internal/domain"
	storage_postgres_gorm "github.com/egorkto/Chat-go/internal/storage/postgres/gorm"
	"gorm.io/gorm"
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
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.User{}, fmt.Errorf(
				"duplicated user, %s: %w",
				err.Error(),
				domain.ErrConflict,
			)
		} else if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			return domain.User{}, fmt.Errorf(
				"check violated, %s: %w",
				err.Error(),
				domain.ErrInvalidArgument,
			)
		}
		return domain.User{}, fmt.Errorf("creating new user: %w", err)
	}
	domainUser := model.ToDomain()

	return domainUser, nil
}
