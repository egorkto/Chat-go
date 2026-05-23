package auth_jwt_transport_http

import (
	"fmt"
	"net/http"

	"github.com/egorkto/Chat-go/internal/domain"
	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
	transport_http_echo "github.com/egorkto/Chat-go/internal/transport/http/echo"
	"github.com/labstack/echo/v5"
)

// SignUpUser godoc
// @Summary      Регистрация нового пользователя
// @Description  Создает новую учетную запись пользователя
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request   body      SignUpRequest  true  "Данные регистрации"
// @Success      201  {object}  AuthResponse "Успешная регистрация"
// @Failure      400  {object}  transport_http.ErrorResponse "Неверный запрос"
// @Failure      409  {object}  transport_http.ErrorResponse "Пользователь уже существует"
// @Failure 	 500  {object}  transport_http.ErrorResponse "Ошибка сервера"
// @Router       /sign-up [post]
func (h *HTTPHandler) SignUp(c *echo.Context) error {
	var request SignUpRequest
	if err := c.Bind(&request); err != nil {
		return transport_http_echo.JSON_Error(
			c,
			"failed to read request body",
			fmt.Errorf("bind request: %w", err),
		)
	}

	if err := c.Validate(request); err != nil {
		return transport_http_echo.JSON_Error(
			c,
			"failed to validate request",
			fmt.Errorf("validate request, %v: %w", request, err),
		)
	}

	domainUser := domain.NewUninitializedUser(request.FullName, request.Login)

	registeredUser, token, err := h.service.SignUp(
		c.Request().Context(),
		domainUser,
		request.Password,
	)
	if err != nil {
		return transport_http_echo.JSON_Error(
			c,
			"failed to sign up",
			fmt.Errorf("sign up user: %w", err),
		)
	}

	access := token.Access
	refresh := token.Refresh

	refreshExpires, err := h.service.GetTokenExpires(refresh)
	if err != nil {
		return transport_http_echo.JSON_Error(
			c,
			"failed to get refresh token expires",
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

	response := responseFromDomain(registeredUser, access)

	return c.JSON(http.StatusCreated, response)
}
