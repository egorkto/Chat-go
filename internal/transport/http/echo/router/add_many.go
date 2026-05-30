package transport_http_echo_router

import (
	"github.com/labstack/echo/v5"
)

func addMany(e *echo.Echo, routes []echo.Route) {
	for _, route := range routes {
		e.AddRoute(route)
	}
}
