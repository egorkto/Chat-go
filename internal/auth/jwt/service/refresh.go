package auth_jwt_service

import (
	"context"
	"fmt"

	auth_jwt_token_manager "github.com/egorkto/Chat-go/internal/auth/jwt/token_manager"
)

func (s *AuthService) Refresh(
	ctx context.Context,
	refreshToken string,
) (auth_jwt_token_manager.JWTPair, error) {
	token, err := s.tokenManager.Verify(refreshToken)
	if err != nil {
		return auth_jwt_token_manager.JWTPair{}, fmt.Errorf(
			"verify refresh token: %w",
			err,
		)
	}

	user, err := s.storage.GetUserByID(ctx, token.UserID)
	if err != nil {
		return auth_jwt_token_manager.JWTPair{}, fmt.Errorf(
			"gete user by id: %w",
			err,
		)
	}

	pair, err := s.tokenManager.Generate(user.ID(), user.Login())
	if err != nil {
		return auth_jwt_token_manager.JWTPair{}, fmt.Errorf(
			"generate: %w",
			err,
		)
	}

	return pair, nil
}
