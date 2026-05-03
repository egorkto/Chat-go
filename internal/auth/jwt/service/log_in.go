package auth_jwt_service

import (
	"context"
	"fmt"

	"github.com/egorkto/Chat-go/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) LogIn(
	ctx context.Context,
	login string,
	password string,
) (domain.User, domain.JWT, error) {
	user, hashedPassword, err := s.storage.GetUserByLogin(ctx, login)
	if err != nil {
		return domain.User{}, domain.JWT{}, fmt.Errorf("getting user: %w", err)
	}

	if err := compare(password, hashedPassword); err != nil {
		return domain.User{}, domain.JWT{}, fmt.Errorf("comparing password: %w", err)
	}

	jwt, err := s.tokenManager.Generate(user)
	if err != nil {
		return domain.User{}, domain.JWT{}, fmt.Errorf("generating jwt token: %w", err)
	}

	return user, jwt, nil
}

func compare(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
