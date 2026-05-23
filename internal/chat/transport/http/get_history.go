package chat_transport_http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/egorkto/Chat-go/internal/domain"
	transport_http_echo "github.com/egorkto/Chat-go/internal/transport/http/echo"
	"github.com/labstack/echo/v5"
)

const (
	defaultLimit  = 20
	defaultOffset = 0
)

func (h HTTPHandler) GetMessages(c *echo.Context) error {
	limit := defaultLimit
	offset := defaultOffset

	var err error

	limitParam := c.QueryParam("limit")
	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			return transport_http_echo.JSON_Error(
				c,
				"failed to get history",
				fmt.Errorf("atoi limit: %w", domain.ErrInvalidArgument),
			)
		}
	}

	offsetParam := c.QueryParam("offset")
	if offsetParam != "" {
		offset, err = strconv.Atoi(offsetParam)
		if err != nil {
			return transport_http_echo.JSON_Error(
				c,
				"failed to get history",
				fmt.Errorf("atoi offset: %w", err),
			)
		}
	}

	messages, err := h.service.GetMessages(c.Request().Context(), &limit, &offset)
	if err != nil {
		return transport_http_echo.JSON_Error(
			c,
			"failed to get history",
			fmt.Errorf("get messages: %w", err),
		)
	}

	response := dtoFromDomains(messages)

	c.JSON(http.StatusOK, response)

	return nil
}
