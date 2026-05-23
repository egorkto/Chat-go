package users_transport_http

import (
	"fmt"
	"net/http"
	"strconv"

	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
	transport_http_echo "github.com/egorkto/Chat-go/internal/transport/http/echo"
	"github.com/labstack/echo/v5"
)

// GetUser godoc
// @Summary      Получение пользователя
// @Description  Возвращает данные существующего пользователя
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id        path      int true                     "Идентификатор пользователя"
// @Success      200       {object}  UserDTOResponse              "Данные пользователя"
// @Failure      401       {object}  transport_http.ErrorResponse "Неавторизованный запрос"
// @Failure      404       {object}  transport_http.ErrorResponse "Пользователь не найден"
// @Failure 	 500       {object}  transport_http.ErrorResponse "Ошибка сервера"
// @Router       /users/{id} [get]
func (h *HTTPHandler) GetUser(c *echo.Context) error {
	idParam := c.Param("id")
	if idParam == "" {
		return c.JSON(
			http.StatusBadRequest,
			transport_http.ErrorResponse{
				Message: "Bad Request",
				Err:     "User ID parameter is missing",
			},
		)
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			transport_http.ErrorResponse{
				Message: "Bad Request",
				Err:     "User ID parameter must be an integer",
			},
		)
	}

	user, err := h.service.GetUser(c.Request().Context(), id)
	if err != nil {
		return transport_http_echo.JSON_Error(
			c,
			"Unauthorized",
			fmt.Errorf("get user: %w", err),
		)
	}

	response := domainToDTO(user)

	return c.JSON(http.StatusOK, response)
}
