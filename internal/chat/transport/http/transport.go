package chat_transport_http

import (
	"context"
	"net/http"

	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/gorilla/websocket"
)

type HTTPHandler struct {
	service ChatService
	hub     WSHub
}

type ChatService interface {
	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)
	GetMessages(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.Message, error)
}

type WSHub interface {
	Upgrade(rw http.ResponseWriter, r *http.Request) (*websocket.Conn, error)
	ReadFrom(conn *websocket.Conn, user domain.User)
}

func New(service ChatService, hub WSHub) *HTTPHandler {
	return &HTTPHandler{
		service: service,
		hub:     hub,
	}
}
