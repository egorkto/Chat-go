package auth_jwt_service

import (
	"context"

	auth_jwt_token_manager "github.com/egorkto/Chat-go/internal/auth/jwt/token_manager"
	"github.com/egorkto/Chat-go/internal/domain"
)

type AuthService struct {
	tokenManager TokenManager
	storage      UsersStorage
	validator    Validator
}

type TokenManager interface {
	Generate(userID int, userLogin string) (auth_jwt_token_manager.JWTPair, error)
	Verify(token string) (auth_jwt_token_manager.Token, error)
}

type Validator interface {
	ValidatePassword(pass string, user domain.User) error
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
