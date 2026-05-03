package auth_jwt_service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/egorkto/Chat-go/internal/domain"
)

func (s *AuthService) Refresh(
	ctx context.Context,
	refreshToken string,
) (domain.JWT, error) {
	token, err := s.tokenManager.Verify(refreshToken)
	if err != nil {
		return domain.JWT{}, fmt.Errorf(
			"verifying refresh token: %w",
			err,
		)
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return domain.JWT{}, fmt.Errorf(
			"getting subject from token: %w",
			err,
		)
	}

	splited := strings.Split(sub, ":")
	if len(splited) != 2 {
		return domain.JWT{}, fmt.Errorf(
			"invalid subject format: %s",
			sub,
		)
	}

	idStr := splited[1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return domain.JWT{}, fmt.Errorf(
			"converting id to integer: %w",
			err,
		)
	}

	user, err := s.storage.GetUserByID(ctx, id)
	if err != nil {
		return domain.JWT{}, fmt.Errorf(
			"getting user by id: %w",
			err,
		)
	}

	jwt, err := s.tokenManager.Generate(user)
	if err != nil {
		return domain.JWT{}, fmt.Errorf(
			"generating jwt token: %w",
			err,
		)
	}

	return jwt, nil
}
