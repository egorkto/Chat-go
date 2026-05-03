package users_service

import (
	"context"

	"github.com/egorkto/Chat-go/internal/domain"
)

type UsersService struct {
	storage UsersStorage
}

type UsersStorage interface {
	GetUserByID(ctx context.Context, id int) (domain.User, error)
}

func New(s UsersStorage) *UsersService {
	return &UsersService{
		storage: s,
	}
}
