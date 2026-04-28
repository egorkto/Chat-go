package auth_service

import (
	"context"
	"fmt"

	"github.com/egorkto/Chat-go/internal/auth"
	"github.com/egorkto/Chat-go/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) SignUp(
	ctx context.Context,
	user domain.User,
	creds auth.Credentials,
) (domain.User, domain.JWT, error) {
	if err := user.Validate(); err != nil {
		return domain.User{}, domain.JWT{}, fmt.Errorf("validating user: %w", err)
	}

	if err := s.validator.ValidatePassword(creds.Password, user.FullName()); err != nil {
		return domain.User{}, domain.JWT{}, fmt.Errorf("validating password: %w", err)
	}

	hashedPass, err := hashed(creds.Password)
	if err != nil {
		return domain.User{}, domain.JWT{}, fmt.Errorf("hashing password: %w", err)
	}

	createdUser, err := s.storage.CreateUser(ctx, user, hashedPass)
	if err != nil {
		return domain.User{}, domain.JWT{}, fmt.Errorf("creating user: %w", err)
	}

	jwt, err := s.generator.Generate(createdUser)
	if err != nil {
		return domain.User{}, domain.JWT{}, fmt.Errorf("generating jwt token: %w", err)
	}

	return createdUser, jwt, nil
}

func hashed(pass string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate from password: %w", err)
	}

	return string(hashed), nil
}
