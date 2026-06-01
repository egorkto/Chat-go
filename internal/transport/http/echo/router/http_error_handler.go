package transport_http_echo_router

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/egorkto/Chat-go/internal/domain"
	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
	"github.com/labstack/echo/v5"
)

func HTTPErrorHandler(c *echo.Context, hErr error) {
	if hErr == nil {
		return
	}

	res, errUnwrap := echo.UnwrapResponse(c.Response())
	if errUnwrap != nil {
		c.Logger().Error("unwrap response", slog.String("err", errUnwrap.Error()))
		return
	} else if res != nil && res.Committed {
		c.Logger().Error("response already commited")
		return
	}

	var errResponse transport_http.ErrorResponse
	var code int

	var coder echo.HTTPStatusCoder
	if errors.As(hErr, &coder) {
		code = coder.StatusCode()
		errResponse.Message = http.StatusText(code)
	} else {
		code = transport_http.ErrorToHTTPCode(hErr)
		errResponse.Message = transport_http.GetHumanized(hErr)
	}

	var response interface{}

	var valErr domain.ValidationError
	if errors.As(hErr, &valErr) {
		response = transport_http.ValidationErrorResponse{
			ErrorResponse: errResponse,
			Details:       valErr.Errors,
		}
	}

	if code >= 400 && code < 500 {
		c.Logger().Warn("client error", slog.String("err", hErr.Error()))
	} else {
		c.Logger().Error("internal error", slog.String("err", hErr.Error()))
	}

	if err := c.JSON(code, response); err != nil {
		c.Logger().Error("json response", slog.String("err", err.Error()))
		return
	}
}
