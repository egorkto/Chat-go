package auth_jwt_transport_http

import (
	"log/slog"
	"net/http"
	"time"

	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
	"github.com/labstack/echo/v5"
)

func (h *HTTPHandler) Refresh(e *echo.Context) error {
	refreshToken, err := e.Cookie("refresh_token")
	if err != nil {
		return e.JSON(
			http.StatusUnauthorized,
			transport_http.ErrorResponse{
				Message: "refresh token is missing",
				Err:     err.Error(),
			})
	}

	token, err := h.service.Refresh(
		e.Request().Context(),
		refreshToken.Value,
	)
	if err != nil {
		e.Logger().Error("refresh token", slog.String("err", err.Error()))
		code := transport_http.ErrorToHTTPCode(err)
		return e.JSON(
			code,
			transport_http.ErrorResponse{
				Message: "failed to refresh token",
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

	return e.JSON(http.StatusCreated, token.Access)
}
