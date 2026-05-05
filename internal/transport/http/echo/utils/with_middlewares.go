package transport_http_echo_utils

import "github.com/labstack/echo/v5"

func WithMiddlewares(routes []echo.Route, middlewares ...echo.MiddlewareFunc) {
	for i := range routes {
		routes[i].Middlewares = append(routes[i].Middlewares, middlewares...)
	}
}
