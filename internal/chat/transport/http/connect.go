package chat_transport_http

import (
	"fmt"

	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/labstack/echo/v5"
)

func (h HTTPHandler) Connect(c *echo.Context) error {
	userID, ok := c.Get("id").(int)
	if !ok {
		return fmt.Errorf("get id: %w", domain.ErrUnauthorized)
	}

	user, err := h.service.GetUser(c.Request().Context(), userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	conn, err := h.hub.Upgrade(c.Response(), c.Request())
	if err != nil {
		return fmt.Errorf("upgrade connection: %w", err)
	}
	defer conn.Close()

	h.hub.ReadFrom(conn, user)

	return nil
}
