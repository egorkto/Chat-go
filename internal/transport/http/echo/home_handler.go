package transport_http_echo

import (
	"fmt"

	"github.com/labstack/echo/v5"
)

func HomeHandler(c *echo.Context) error {
	err := c.File("web/index.html")
	if err != nil {
		return JSON_Error(
			c,
			"failed to load home page",
			fmt.Errorf("echo context file: %w", err),
		)
	}
	return nil
}
