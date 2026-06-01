package auth_jwt_transport_http

import (
	"fmt"
	"net/http"

	"github.com/egorkto/Chat-go/internal/domain"
	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
	"github.com/labstack/echo/v5"
)

// SignUpUser godoc
// @Summary      Регистрация нового пользователя
// @Description  Создает новую учетную запись пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request   body      SignUpRequest  true  "Данные регистрации"
// @Success      201  {object}  AuthResponse "Успешная регистрация"
// @Failure      400  {object}  ValidationErrorResponse "Неверный запрос"
// @Failure      409  {object}  ErrorResponse "Пользователь уже существует"
// @Failure 	 500  {object}  ErrorResponse "Ошибка сервера"
// @Router       /sign-up [post]
func (h *HTTPHandler) SignUp(c *echo.Context) error {
	var request SignUpRequest
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

	domainUser := domain.NewUninitializedUser(request.FullName, request.Login)

	registeredUser, pair, err := h.service.SignUp(
		c.Request().Context(),
		domainUser,
		request.Password,
	)
	if err != nil {
		return fmt.Errorf("sign up: %w", err)
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

	response := responseFromDomain(registeredUser, access.Signed)

	return c.JSON(http.StatusCreated, response)
}
