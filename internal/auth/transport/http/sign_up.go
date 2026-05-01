package auth_transport

import (
	"log/slog"
	"net/http"

	"github.com/egorkto/Chat-go/http_response"
	"github.com/egorkto/Chat-go/internal/auth"
	"github.com/egorkto/Chat-go/internal/domain"
	echo_dto "github.com/egorkto/Chat-go/internal/echo"
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
			echo_dto.ErrorResponse{
				Message: "failed to read request body",
				Err:     err.Error(),
			})
	}

	if err := e.Validate(request); err != nil {
		e.Logger().Error("validate request", slog.String("err", err.Error()))
		return e.JSON(
			http.StatusBadRequest,
			echo_dto.ErrorResponse{
				Message: "failed to validate request",
				Err:     err.Error(),
			})
	}

	domainUser := domain.NewUninitializedUser(request.FullName)
	credentials := auth.NewCredentials(request.Password)

	registeredUser, token, err := h.service.SignUp(
		e.Request().Context(),
		domainUser,
		credentials,
	)
	if err != nil {
		e.Logger().Error("sign up user", slog.String("err", err.Error()))
		code := http_response.ErrorToHTTPCode(err)
		return e.JSON(code, echo_dto.ErrorResponse{
			Message: "failed to sign up",
			Err:     err.Error(),
		})
	}

	response := responseFromDomain(registeredUser, token)

	return e.JSON(http.StatusCreated, response)
}

func responseFromDomain(u domain.User, t domain.JWT) SignUpResponse {
	return SignUpResponse{
		User: UserDTO{
			u.ID(),
			u.Version(),
			u.FullName(),
		},
		AccessToken:  t.Access,
		RefreshToken: t.Refresh,
	}
}
