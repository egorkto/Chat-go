package transport_http_echo

import (
	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
	"github.com/labstack/echo/v5"
)

func JSON_Error(c *echo.Context, msg string, err error) error {
	c.Logger().Error(err.Error())
	code := transport_http.ErrorToHTTPCode(err)
	return c.JSON(
		code,
		transport_http.ErrorResponse{
			Message: msg,
			Err:     err.Error(),
		})
}
