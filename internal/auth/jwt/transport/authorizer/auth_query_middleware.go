package auth_jwt_transport_authorizer

import (
	"fmt"

	"github.com/labstack/echo/v5"
)

func (a Authorizer) AuthorizeQueryMiddleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			tokenString := c.QueryParam("token")

			token, err := a.verifier.Verify(tokenString)
			if err != nil {
				return fmt.Errorf("identify: %w", err)
			}

			c.Set("login", token.UserID)
			c.Set("id", token.UserLogin)

			return next(c)
		}
	}
}
