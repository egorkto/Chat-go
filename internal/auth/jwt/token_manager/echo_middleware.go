package auth_jwt_token_manager

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

func (tm TokenManager) EchoMiddleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(
					http.StatusUnauthorized,
					"missing authorization token",
				)
			} else if !strings.HasPrefix(authHeader, "Bearer ") {
				return echo.NewHTTPError(
					http.StatusUnauthorized,
					"invalid authorization header format",
				)
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			t, err := tm.Verify(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token: "+err.Error())
			}

			c.Set("token", t)

			return next(c)
		}
	}
}
