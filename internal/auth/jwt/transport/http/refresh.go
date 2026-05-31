package auth_jwt_transport_http

import (
	"fmt"
	"net/http"

	"github.com/egorkto/Chat-go/internal/domain"
	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
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
		return fmt.Errorf("get refresh token cookie, %s: %w", err.Error(), domain.ErrUnauthorized)
	}

	pair, err := h.service.Refresh(
		c.Request().Context(),
		refreshToken.Value,
	)
	if err != nil {
		return fmt.Errorf("refresh tokenЖ %w", err)
	}

	access := pair.Access
	refresh := pair.Refresh

	cookie := transport_http.NewCookie(
		"refresh_token",
		refresh.Signed,
		refresh.ExpiredAt,
		"/refresh",
		true,
	)

	c.SetCookie(cookie)

	return c.JSON(http.StatusCreated, access)
}
