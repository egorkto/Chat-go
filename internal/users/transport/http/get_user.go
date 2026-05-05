package users_transport_http

import (
	"net/http"
	"strconv"

	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
	"github.com/labstack/echo/v5"
)

func (h *HTTPHandler) GetUser(e *echo.Context) error {
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

	user, err := h.service.GetUser(e.Request().Context(), id)
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
