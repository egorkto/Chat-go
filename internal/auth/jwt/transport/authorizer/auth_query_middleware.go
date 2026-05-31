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
				return fmt.Errorf("verify: %w", err)
			}

			c.Set("id", token.UserID)
			c.Set("login", token.UserLogin)

			return next(c)
		}
	}
}
