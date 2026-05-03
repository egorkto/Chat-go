package auth_jwt_service

import (
	"context"

	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	tokenManager TokenManager
	storage      UsersStorage
	validator    Validator
}

type TokenManager interface {
	Generate(u domain.User) (domain.JWT, error)
	Verify(token string) (*jwt.Token, error)
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

	GetUserByLogin(
		ctx context.Context,
		login string,
	) (domain.User, string, error)

	GetUserByID(
		ctx context.Context,
		id int,
	) (domain.User, error)
}

func New(g TokenManager, s UsersStorage, v Validator) *AuthService {
	return &AuthService{
		tokenManager: g,
		storage:      s,
		validator:    v,
	}
}
