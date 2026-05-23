package auth_jwt_transport_http

import (
	"fmt"
	"net/http"

	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
	transport_http_echo "github.com/egorkto/Chat-go/internal/transport/http/echo"
	"github.com/labstack/echo/v5"
)

// Refresh godoc
// @Summary      Обновление токена
// @Description  Возвращает новые access и refresh токены
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200       {object}  string                       "Обновленный access токен"
// @Failure      401       {object}  transport_http.ErrorResponse "Неавторизованный запрос"
// @Failure 	 500       {object}  transport_http.ErrorResponse "Ошибка сервера"
// @Router       /refresh [post]
func (h *HTTPHandler) Refresh(c *echo.Context) error {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		return transport_http_echo.JSON_Error(
			c,
			"Unauthorized",
			fmt.Errorf("get refresh token cookie: %w", err),
		)
	}

	domainJWT, err := h.service.Refresh(
		c.Request().Context(),
		refreshToken.Value,
	)
	if err != nil {
		return transport_http_echo.JSON_Error(
			c,
			"Unauthorized",
			fmt.Errorf("refresh token"),
		)
	}

	access := domainJWT.Access
	refresh := domainJWT.Refresh

	refreshExpires, err := h.service.GetTokenExpires(refresh)
	if err != nil {
		return transport_http_echo.JSON_Error(
			c,
			"Unauthorized",
			fmt.Errorf("get refresh token expires: %w", err),
		)
	}

	cookie := transport_http.NewCookie(
		"refresh_token",
		refresh,
		refreshExpires,
		"/refresh",
		true,
	)

	c.SetCookie(cookie)

	return c.JSON(http.StatusCreated, access)
}
