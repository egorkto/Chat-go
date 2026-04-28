package auth_service

import (
	"context"

	"github.com/egorkto/Chat-go/internal/domain"
)

type AuthService struct {
	generator JWTGenerator
	storage   UsersStorage
	validator Validator
}

type JWTGenerator interface {
	Generate(u domain.User) (domain.JWT, error)
}

type Validator interface {
	ValidatePassword(pass string, name string) error
}

type UsersStorage interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
		password string,
	) (domain.User, error)
}

func New(g JWTGenerator, s UsersStorage, v Validator) *AuthService {
	return &AuthService{
		generator: g,
		storage:   s,
		validator: v,
	}
}
