package users_transport_http

import (
	"context"

	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type HTTPHandler struct {
	service UsersService
}

type UsersService interface {
	GetUser(ctx context.Context, id int, token *jwt.Token) (domain.User, error)
}

func New(s UsersService) *HTTPHandler {
	return &HTTPHandler{
		service: s,
	}
}

func (h *HTTPHandler) Routes() []echo.Route {
	return []echo.Route{
		{
			Method:  "GET",
			Path:    "/users/:id",
			Handler: h.GetUser,
		},
	}
}
