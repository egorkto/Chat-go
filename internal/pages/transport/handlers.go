package pages_transport

import (
	"fmt"

	"github.com/labstack/echo/v5"
)

func (h HTTPHandler) Home(c *echo.Context) error {
	err := c.File("web/index.html")
	if err != nil {
		return fmt.Errorf("file: %w", err)
	}
	return nil
}

func (h HTTPHandler) SignUp(c *echo.Context) error {
	err := c.File("web/sign-up.html")
	if err != nil {
		return fmt.Errorf("file: %w", err)
	}
	return nil
}

func (h HTTPHandler) LogIn(c *echo.Context) error {
	err := c.File("web/log-in.html")
	if err != nil {
		return fmt.Errorf("file: %w", err)
	}
	return nil
}
