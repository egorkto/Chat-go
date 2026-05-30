package pages_transport

import (
	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
	"github.com/labstack/echo/v5"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

func (h HTTPHandler) Routes(mc transport_http.MiddlewaresContainer) []echo.Route {
	return []echo.Route{
		{
			Method:  "GET",
			Path:    "/",
			Handler: h.Home,
		},
		{
			Method:  "GET",
			Path:    "/sign-up",
			Handler: h.SignUp,
		},
		{
			Method:  "GET",
			Path:    "/log-in",
			Handler: h.LogIn,
		},
		{
			Method:  "GET",
			Path:    "/swagger/*",
			Handler: echoSwagger.WrapHandler,
		},
	}
}
