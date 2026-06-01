package auth_jwt_transport_http

import (
	"fmt"
	"net/http"

	"github.com/egorkto/Chat-go/internal/domain"
	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
	"github.com/labstack/echo/v5"
)

// LogInUser godoc
// @Summary      Авторизация пользователя
// @Description  Авторизирует существующего пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request   body      LogInRequest  true           "Данные авторизации"
// @Success      200       {object}  AuthResponse                "Успешная авторизация"
// @Failure      400       {object}  ValidationErrorResponse "Неверный запрос"
// @Failure		 404       {object}  ErrorResponse "Пользователь не найден"
// @Failure 	 500       {object}  ErrorResponse "Ошибка сервера"
// @Router       /log-in [post]
func (h *HTTPHandler) LogIn(c *echo.Context) error {
	var request LogInRequest
	if err := c.Bind(&request); err != nil {
		return fmt.Errorf(
			"bind request, %s: %w",
			err.Error(),
			domain.NewValidationError(map[string]string{
				"body": "invalid JSON format",
			}))
	}

	if err := c.Validate(request); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	user, pair, err := h.service.LogIn(
		c.Request().Context(),
		request.Login,
		request.Password,
	)
	if err != nil {
		return fmt.Errorf("log in: %w", err)
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

	response := responseFromDomain(user, access.Signed)

	return c.JSON(http.StatusOK, response)
}
