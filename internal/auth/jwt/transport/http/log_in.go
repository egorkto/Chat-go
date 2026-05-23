package auth_jwt_transport_http

import (
	"fmt"
	"net/http"

	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
	transport_http_echo "github.com/egorkto/Chat-go/internal/transport/http/echo"
	"github.com/labstack/echo/v5"
)

// LogInUser godoc
// @Summary      Авторизация пользователя
// @Description  Авторизирует существующего пользователя
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request   body      LogInRequest  true           "Данные авторизации"
// @Success      200       {object}  AuthResponse                "Успешная авторизация"
// @Failure      400       {object}  transport_http.ErrorResponse "Неверный запрос"
// @Failure		 404       {object}  transport_http.ErrorResponse "Пользователь не найден"
// @Failure 	 500       {object}  transport_http.ErrorResponse "Ошибка сервера"
// @Router       /log-in [post]
func (h *HTTPHandler) LogIn(c *echo.Context) error {
	var request LogInRequest
	if err := c.Bind(&request); err != nil {
		return transport_http_echo.JSON_Error(
			c,
			"failed to log in",
			fmt.Errorf("bind request: %w", err),
		)
	}

	if err := c.Validate(request); err != nil {
		return transport_http_echo.JSON_Error(
			c,
			"wrong login or password",
			fmt.Errorf("validate request: %w", err),
		)
	}

	user, domainJWT, err := h.service.LogIn(
		c.Request().Context(),
		request.Login,
		request.Password,
	)
	if err != nil {
		return transport_http_echo.JSON_Error(
			c,
			"wrong login or password",
			fmt.Errorf("log in user: %w", err),
		)
	}

	access := domainJWT.Access
	refresh := domainJWT.Refresh

	refreshExpires, err := h.service.GetTokenExpires(refresh)
	if err != nil {
		return transport_http_echo.JSON_Error(
			c,
			"failed to log in",
			fmt.Errorf("get refreshg token expires: %w", err),
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

	response := responseFromDomain(user, access)

	return c.JSON(http.StatusOK, response)
}
