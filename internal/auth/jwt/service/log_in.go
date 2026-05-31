package auth_jwt_service

import (
	"context"
	"fmt"

	auth_jwt_token_manager "github.com/egorkto/Chat-go/internal/auth/jwt/token_manager"
	"github.com/egorkto/Chat-go/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) LogIn(
	ctx context.Context,
	login string,
	password string,
) (domain.User, auth_jwt_token_manager.JWTPair, error) {
	user, hashedPassword, err := s.storage.GetUserByLogin(ctx, login)
	if err != nil {
		return domain.User{}, auth_jwt_token_manager.JWTPair{}, fmt.Errorf(
			"get user by login: %w",
			err,
		)
	}

	if err := compare(password, hashedPassword); err != nil {
		return domain.User{}, auth_jwt_token_manager.JWTPair{}, fmt.Errorf(
			"compare password, %s: %w",
			err.Error(),
			domain.ErrPasswordNotMatch,
		)
	}

	pair, err := s.tokenManager.Generate(user.ID(), user.Login())
	if err != nil {
		return domain.User{}, auth_jwt_token_manager.JWTPair{}, fmt.Errorf(
			"generate jwt token: %w",
			err,
		)
	}

	return user, pair, nil
}

func compare(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
