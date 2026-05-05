package auth_jwt_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/labstack/echo/v5"
)

type HTTPHandler struct {
	service AuthService
}

type AuthService interface {
	SignUp(
		ctx context.Context,
		user domain.User,
		password string,
	) (domain.User, domain.JWT, error)

	LogIn(
		ctx context.Context,
		login string,
		password string,
	) (domain.User, domain.JWT, error)

	Refresh(
		ctx context.Context,
		refreshToken string,
	) (domain.JWT, error)

	GetTokenExpires(token string) (time.Time, error)
}

func New(s AuthService) *HTTPHandler {
	return &HTTPHandler{
		service: s,
	}
}

func (h *HTTPHandler) Routes() []echo.Route {
	return []echo.Route{
		{
			Method:  http.MethodPost,
			Path:    "/sign-up",
			Handler: h.SignUp,
		},
		{
			Method:  http.MethodPost,
			Path:    "/log-in",
			Handler: h.LogIn,
		},
		{
			Method:  http.MethodPost,
			Path:    "/refresh",
			Handler: h.Refresh,
		},
	}
}
