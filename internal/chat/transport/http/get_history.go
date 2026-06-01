package chat_transport_http

import (
	"fmt"
	"net/http"
	"strconv"

	chat_transport "github.com/egorkto/Chat-go/internal/chat/transport"
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/labstack/echo/v5"
)

const (
	defaultLimit  = 20
	defaultOffset = 0
)

// GetHistory godoc
// @Summary      Получение истории чата
// @Description  Возвращает историю сообщений в чате
// @Tags         chat
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit     query     int false       "Лимит сообщений"
// @Param        offset    query     int false       "Смещение"
// @Success      200       {object}  GetHistoryResponse "Данные пользователя"
// @Failure      401       {object}  ErrorResponse   "Неавторизованный запрос"
// @Failure 	 500       {object}  ErrorResponse   "Ошибка сервера"
// @Router       /chat/history [get]
func (h HTTPHandler) GetHistory(c *echo.Context) error {
	limit := defaultLimit
	offset := defaultOffset

	var err error

	limitParam := c.QueryParam("limit")
	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			return fmt.Errorf("atoi limit: %w", domain.NewValidationError(map[string]string{
				"limit": "limit query param is not a number",
			}))
		}
	}

	offsetParam := c.QueryParam("offset")
	if offsetParam != "" {
		offset, err = strconv.Atoi(offsetParam)
		if err != nil {
			return fmt.Errorf("atoi offset: %w", domain.NewValidationError(map[string]string{
				"offset": "offset query param is not a number",
			}))
		}
	}

	messages, err := h.service.GetMessages(c.Request().Context(), &limit, &offset)
	if err != nil {
		return fmt.Errorf("get messages: %w", err)
	}

	response := chat_transport.DtoFromDomains(messages)

	c.JSON(http.StatusOK, response)

	return nil
}
