package auth_jwt_transport_http

import (
	"log/slog"
	"net/http"

	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
	"github.com/labstack/echo/v5"
)

func (h *HTTPHandler) LogIn(e *echo.Context) error {
	var request LogInRequest
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

	user, domainJWT, err := h.service.LogIn(
		e.Request().Context(),
		request.Login,
		request.Password,
	)
	if err != nil {
		e.Logger().Error("log in user", slog.String("err", err.Error()))
		code := transport_http.ErrorToHTTPCode(err)
		return e.JSON(
			code,
			transport_http.ErrorResponse{
				Message: "failed to log in",
				Err:     err.Error(),
			})
	}

	access := domainJWT.Access
	refresh := domainJWT.Refresh

	refreshExpires, err := h.service.GetTokenExpires(refresh)
	if err != nil {
		e.Logger().Error("get refresh token expires", slog.String("err", err.Error()))
		code := transport_http.ErrorToHTTPCode(err)
		return e.JSON(
			code,
			transport_http.ErrorResponse{
				Message: "failed to get refresh token expires",
				Err:     err.Error(),
			})
	}

	cookie := transport_http.NewCookie(
		"refresh_token",
		refresh,
		refreshExpires,
		"/refresh",
		true,
	)

	e.SetCookie(cookie)

	response := responseFromDomain(user, access)

	return e.JSON(http.StatusOK, response)
}
