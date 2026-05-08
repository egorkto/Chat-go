package users_storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/egorkto/Chat-go/internal/domain"
	storage_postgres "github.com/egorkto/Chat-go/internal/storage/postgres"
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
		fmt.Println(err.Error())
		fmt.Println(errors.Is(err, storage_postgres.ErrDuplicatedKey))
		if errors.Is(err, storage_postgres.ErrDuplicatedKey) {
			return domain.User{}, fmt.Errorf(
				"duplicated user, %s: %w",
				err.Error(),
				domain.ErrConflict,
			)
		}
		return domain.User{}, fmt.Errorf("creating new user: %w", err)
	}
	domainUser := model.ToDomain()

	return domainUser, nil
}
