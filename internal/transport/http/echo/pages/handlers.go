package transport_http_echo_pages

import (
	"fmt"

	transport_http_echo "github.com/egorkto/Chat-go/internal/transport/http/echo"
	"github.com/labstack/echo/v5"
)

func Home(c *echo.Context) error {
	err := c.File("web/index.html")
	if err != nil {
		return transport_http_echo.JSON_Error(
			c,
			"failed to load home page",
			fmt.Errorf("echo context file: %w", err),
		)
	}
	return nil
}

func SignUp(c *echo.Context) error {
	err := c.File("web/sign-up.html")
	if err != nil {
		return transport_http_echo.JSON_Error(
			c,
			"failed to load home page",
			fmt.Errorf("echo context file: %w", err),
		)
	}
	return nil
}

func LogIn(c *echo.Context) error {
	err := c.File("web/log-in.html")
	if err != nil {
		return transport_http_echo.JSON_Error(
			c,
			"failed to load home page",
			fmt.Errorf("echo context file: %w", err),
		)
	}
	return nil
}
