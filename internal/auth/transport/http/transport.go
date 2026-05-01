package auth_transport

import (
	"context"
	"net/http"

	"github.com/egorkto/Chat-go/internal/auth"
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/labstack/echo/v5"
)

type HTTPHandler struct {
	service AuthService
}

func New(s AuthService) *HTTPHandler {
	return &HTTPHandler{
		service: s,
	}
}

type AuthService interface {
	SignUp(
		ctx context.Context,
		user domain.User,
		creds auth.Credentials,
	) (domain.User, domain.JWT, error)
}

func (h *HTTPHandler) Routes() []echo.Route {
	return []echo.Route{
		{
			Method:  http.MethodPost,
			Path:    "/sign-up",
			Handler: h.SignUp,
		},
	}
}
