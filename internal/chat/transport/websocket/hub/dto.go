package chat_transport_websocket_hub

import (
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	User domain.User
}

type MessageInput struct {
	Text string `json:"text"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
