package auth_service

import (
	"context"

	"github.com/egorkto/Chat-go/internal/auth"
	"github.com/egorkto/Chat-go/internal/domain"
)

type AuthService struct {
	generator auth.JWTGenerator
	storage   UsersStorage
}

type UsersStorage interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
		password string,
	) (domain.User, error)
}

func New(g auth.JWTGenerator, storage UsersStorage) *AuthService {
	return &AuthService{
		generator: g,
		storage:   storage,
	}
}
