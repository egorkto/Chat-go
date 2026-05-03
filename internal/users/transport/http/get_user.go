package users_transport_http

import (
	"net/http"
	"strconv"

	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func (h *HTTPHandler) GetUser(e *echo.Context) error {
	token, ok := e.Get("token").(*jwt.Token)
	if !ok {
		return e.JSON(
			http.StatusUnauthorized,
			transport_http.ErrorResponse{
				Message: "Unauthorized",
				Err:     "Token in echo context is missing",
			},
		)
	}

	idParam := e.Param("id")
	if idParam == "" {
		return e.JSON(
			http.StatusBadRequest,
			transport_http.ErrorResponse{
				Message: "Bad Request",
				Err:     "User ID parameter is missing",
			},
		)
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return e.JSON(
			http.StatusBadRequest,
			transport_http.ErrorResponse{
				Message: "Bad Request",
				Err:     "User ID parameter must be an integer",
			},
		)
	}

	user, err := h.service.GetUser(e.Request().Context(), id, token)
	if err != nil {
		return e.JSON(
			http.StatusUnauthorized,
			transport_http.ErrorResponse{
				Message: "Unauthorized",
				Err:     err.Error(),
			},
		)
	}

	response := domainToDTO(user)

	return e.JSON(http.StatusOK, response)
}
