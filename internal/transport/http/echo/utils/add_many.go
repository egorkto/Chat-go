package transport_http_echo_utils

import (
	"github.com/labstack/echo/v5"
)

func AddMany(e *echo.Echo, routes []echo.Route) {
	for _, route := range routes {
		e.AddRoute(route)
	}
}
