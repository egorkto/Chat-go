package users_transport_http

import (
	"context"

	"github.com/egorkto/Chat-go/internal/domain"
)

type HTTPHandler struct {
	service UsersService
}

type UsersService interface {
	GetUser(ctx context.Context, id int) (domain.User, error)
}

func New(s UsersService) *HTTPHandler {
	return &HTTPHandler{
		service: s,
	}
}
