package chat_transport_http

import (
	"context"
	"net/http"

	chat_transport_websocket_hub "github.com/egorkto/Chat-go/internal/chat/transport/websocket/hub"
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

type HTTPHandler struct {
	service ChatService
	hub     WebsocketHub
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

type WebsocketHub interface {
	Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error)
	ReadFrom(c chat_transport_websocket_hub.Client) error
}

func New(service ChatService, hub WebsocketHub) HTTPHandler {
	return HTTPHandler{
		service: service,
		hub:     hub,
	}
}

func (h *HTTPHandler) Routes() []echo.Route {
	return []echo.Route{
		{
			Method:  "GET",
			Path:    "/history",
			Handler: h.GetMessages,
		},
	}
}

func (h *HTTPHandler) ConnectWebsocketRoute() echo.Route {
	return echo.Route{
		Method:  "GET",
		Path:    "/connect",
		Handler: h.Connect,
	}
}
