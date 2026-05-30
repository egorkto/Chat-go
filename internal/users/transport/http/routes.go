package users_transport_http

import (
	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
	"github.com/labstack/echo/v5"
)

func (h *HTTPHandler) Routes(mc transport_http.MiddlewaresContainer) []echo.Route {
	return []echo.Route{
		{
			Method:  "GET",
			Path:    "/users/:id",
			Handler: h.GetUser,
			Middlewares: []echo.MiddlewareFunc{
				mc.HeaderAuth,
			},
		},
	}
}
