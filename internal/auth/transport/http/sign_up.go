package auth_transport

import (
	"log/slog"
	"net/http"

	"github.com/egorkto/Chat-go/http_response"
	"github.com/egorkto/Chat-go/internal/auth"
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/labstack/echo/v5"
)

func (h *HTTPHandler) SignUp(e *echo.Context) error {
	var request SignUpRequest
	if err := e.Bind(&request); err != nil {
		e.Logger().Error("bind request", slog.String("err", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to read request body")
	}

	if err := e.Validate(request); err != nil {
		e.Logger().Error("validate request", slog.String("err", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to validate request")
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
		return echo.NewHTTPError(code, "failed to sign up")
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
