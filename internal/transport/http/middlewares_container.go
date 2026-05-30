package transport_http

import "github.com/labstack/echo/v5"

type MiddlewaresContainer struct {
	HeaderAuth echo.MiddlewareFunc
	QueryAuth  echo.MiddlewareFunc
}
