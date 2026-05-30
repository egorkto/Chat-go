package auth_jwt_transport_authorizer

import (
	"fmt"
	"strings"

	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/labstack/echo/v5"
)

func (a Authorizer) AuthorizeHeaderMiddleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader, err := getAuthHeader(c)
			if err != nil {
				return fmt.Errorf("validate auth header: %w", err)
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := a.verifier.Verify(tokenString)
			if err != nil {
				return fmt.Errorf("verify: %w", err)
			}

			c.Set("login", token.UserID)
			c.Set("id", token.UserLogin)

			return next(c)
		}
	}
}

func getAuthHeader(c *echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf(
			"missing 'Bearer' prefix: %w",
			domain.ErrUnauthorized,
		)
	}

	return authHeader, nil
}
