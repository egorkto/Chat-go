package auth_jwt_transport_http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/egorkto/Chat-go/internal/domain"
	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
	"github.com/labstack/echo/v5"
)

// SignUpUser godoc
// @Summary      Регистрация нового пользователя
// @Description  Создает новую учетную запись пользователя
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request   body      SignUpRequest  true  "Данные регистрации"
// @Success      201  {object}  SignUpResponse "Успешная регистрация"
// @Failure      400  {object}  echo_dto.ErrorResponse "Неверный запрос"
// @Failure      409  {object}  echo_dto.ErrorResponse "Пользователь уже существует"
// @Router       /sign-up [post]
func (h *HTTPHandler) SignUp(e *echo.Context) error {
	var request SignUpRequest
	if err := e.Bind(&request); err != nil {
		e.Logger().Error("bind request", slog.String("err", err.Error()))
		return e.JSON(
			http.StatusBadRequest,
			transport_http.ErrorResponse{
				Message: "failed to read request body",
				Err:     err.Error(),
			})
	}

	if err := e.Validate(request); err != nil {
		e.Logger().Error("validate request", slog.String("err", err.Error()))
		return e.JSON(
			http.StatusBadRequest,
			transport_http.ErrorResponse{
				Message: "failed to validate request",
				Err:     err.Error(),
			})
	}

	domainUser := domain.NewUninitializedUser(request.FullName, request.Login)

	registeredUser, token, err := h.service.SignUp(
		e.Request().Context(),
		domainUser,
		request.Password,
	)
	if err != nil {
		e.Logger().Error("sign up user", slog.String("err", err.Error()))
		code := transport_http.ErrorToHTTPCode(err)
		return e.JSON(
			code,
			transport_http.ErrorResponse{
				Message: "failed to sign up",
				Err:     err.Error(),
			})
	}

	cookie := transport_http.NewCookie(
		"refresh_token",
		token.Refresh,
		time.Now().Add(token.RefreshExpires),
		"/refresh",
		true,
	)

	e.SetCookie(cookie)

	response := responseFromDomain(registeredUser, token.Access)

	return e.JSON(http.StatusCreated, response)
}
