package auth_jwt_transport_http

import (
	"log/slog"
	"net/http"

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

	domainJWT, err := h.service.Refresh(
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

	return e.JSON(http.StatusCreated, access)
}
