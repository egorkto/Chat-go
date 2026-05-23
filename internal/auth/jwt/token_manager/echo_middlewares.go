package auth_jwt_token_manager

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func (tm TokenManager) HeaderMiddleware() func(next echo.HandlerFunc) echo.HandlerFunc {
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
				return echo.NewHTTPError(
					http.StatusUnauthorized,
					fmt.Sprintf(
						"invalid token: %s",
						err.Error(),
					),
				)
			}

			login, id, err := getLoginIDFromToken(t)
			if err != nil {
				return echo.NewHTTPError(
					http.StatusUnauthorized,
					fmt.Sprintf(
						"get login, id from token: %s",
						err.Error(),
					),
				)
			}

			c.Set("login", login)
			c.Set("id", id)

			return next(c)
		}
	}
}

func (tm TokenManager) QueryMiddleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			token := c.QueryParam("token")
			if token == "" {
				return echo.NewHTTPError(
					http.StatusUnauthorized,
					"missing token in query param",
				)
			}

			t, err := tm.Verify(token)
			if err != nil {
				return echo.NewHTTPError(
					http.StatusUnauthorized,
					fmt.Sprintf(
						"invalid token: %s",
						err.Error(),
					),
				)
			}

			login, id, err := getLoginIDFromToken(t)
			if err != nil {
				return echo.NewHTTPError(
					http.StatusUnauthorized,
					fmt.Sprintf(
						"get login, id from token: %s",
						err.Error(),
					),
				)
			}

			c.Set("login", login)
			c.Set("id", id)

			return next(c)
		}
	}
}

func getLoginIDFromToken(token *jwt.Token) (string, int, error) {
	sub, err := token.Claims.GetSubject()
	if err != nil {
		return "", 0, fmt.Errorf(
			"get subject: %w",
			err,
		)
	}

	splited := strings.Split(sub, ":")
	if len(splited) != 2 {
		return "", 0, fmt.Errorf(
			"wrong subject format: %w",
			domain.ErrUnauthorized,
		)
	}

	login := splited[0]
	strID := splited[1]

	id, err := strconv.Atoi(strID)
	if err != nil {
		return "", 0, fmt.Errorf(
			"id is not a int, %w",
			err,
		)
	}

	return login, id, nil
}
