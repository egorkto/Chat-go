package chat_transport_http

import (
	"fmt"
	"log/slog"

	chat_transport_websocket_hub "github.com/egorkto/Chat-go/internal/chat/transport/websocket/hub"
	"github.com/egorkto/Chat-go/internal/domain"
	transport_http_echo "github.com/egorkto/Chat-go/internal/transport/http/echo"
	"github.com/labstack/echo/v5"
)

func (h HTTPHandler) Connect(c *echo.Context) error {
	conn, err := h.hub.Upgrade(c.Response(), c.Request())
	if err != nil {
		return transport_http_echo.JSON_Error(
			c,
			"Failed to connect",
			fmt.Errorf("upgrade connection: %w", err),
		)
	}
	defer conn.Close()

	userID, ok := c.Get("id").(int)
	if !ok {
		return transport_http_echo.JSON_Error(
			c,
			"Unauthorized",
			fmt.Errorf("get id: %w", domain.ErrUnauthorized),
		)
	}

	user, err := h.service.GetUser(c.Request().Context(), userID)
	if err != nil {
		return transport_http_echo.JSON_Error(
			c,
			"Failed to connect",
			fmt.Errorf("get user: %w", err),
		)
	}

	client := chat_transport_websocket_hub.Client{
		Conn: conn,
		User: user,
	}

	if err := h.hub.ReadFrom(client); err != nil {
		c.Logger().Error("read from client", slog.String("err", err.Error()))
	}

	return nil
}
