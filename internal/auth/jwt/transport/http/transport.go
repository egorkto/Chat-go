package auth_jwt_transport_http

import (
	"context"

	auth_jwt_token_manager "github.com/egorkto/Chat-go/internal/auth/jwt/token_manager"
	"github.com/egorkto/Chat-go/internal/domain"
)

type HTTPHandler struct {
	service AuthService
}

type AuthService interface {
	SignUp(
		ctx context.Context,
		user domain.User,
		password string,
	) (domain.User, auth_jwt_token_manager.JWTPair, error)

	LogIn(
		ctx context.Context,
		login string,
		password string,
	) (domain.User, auth_jwt_token_manager.JWTPair, error)

	Refresh(
		ctx context.Context,
		refreshToken string,
	) (auth_jwt_token_manager.JWTPair, error)
}

func New(s AuthService) *HTTPHandler {
	return &HTTPHandler{
		service: s,
	}
}
