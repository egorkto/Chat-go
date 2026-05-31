package auth_jwt_service

import (
	"context"
	"fmt"

	auth_jwt_token_manager "github.com/egorkto/Chat-go/internal/auth/jwt/token_manager"
	"github.com/egorkto/Chat-go/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) SignUp(
	ctx context.Context,
	user domain.User,
	password string,
) (domain.User, auth_jwt_token_manager.JWTPair, error) {
	if err := user.Validate(); err != nil {
		return domain.User{}, auth_jwt_token_manager.JWTPair{}, fmt.Errorf(
			"validate user: %w",
			err,
		)
	}

	if err := s.validator.ValidatePassword(password, user); err != nil {
		return domain.User{}, auth_jwt_token_manager.JWTPair{}, fmt.Errorf(
			"validate password: %w",
			err,
		)
	}

	hashedPass, err := hashed(password)
	if err != nil {
		return domain.User{}, auth_jwt_token_manager.JWTPair{}, fmt.Errorf(
			"hash password: %w",
			err,
		)
	}

	createdUser, err := s.storage.CreateUser(ctx, user, hashedPass)
	if err != nil {
		return domain.User{}, auth_jwt_token_manager.JWTPair{}, fmt.Errorf(
			"create user: %w",
			err,
		)
	}

	pair, err := s.tokenManager.Generate(createdUser.ID(), createdUser.Login())
	if err != nil {
		return domain.User{}, auth_jwt_token_manager.JWTPair{}, fmt.Errorf(
			"generate: %w",
			err,
		)
	}

	return createdUser, pair, nil
}

func hashed(pass string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate from password: %w", err)
	}

	return string(hashed), nil
}
