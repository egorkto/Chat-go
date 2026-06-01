package users_transport_http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/egorkto/Chat-go/internal/domain"
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
// @Failure      401       {object}  ErrorResponse "Неавторизованный запрос"
// @Failure      404       {object}  ErrorResponse "Пользователь не найден"
// @Failure 	 500       {object}  ErrorResponse "Ошибка сервера"
// @Router       /users/{id} [get]
func (h *HTTPHandler) GetUser(c *echo.Context) error {
	idParam := c.Param("id")
	if idParam == "" {
		return fmt.Errorf("id param is empty: %w", domain.NewValidationError(map[string]string{
			"id": "id route param is empty",
		}))
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return fmt.Errorf("atoi id param: %w", err)
	}

	user, err := h.service.GetUser(c.Request().Context(), id)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	response := domainToDTO(user)

	return c.JSON(http.StatusOK, response)
}
