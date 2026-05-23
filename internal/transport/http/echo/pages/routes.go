package transport_http_echo_pages

import (
	"github.com/labstack/echo/v5"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

func Routes() []echo.Route {
	return []echo.Route{
		{
			Method:  "GET",
			Path:    "/",
			Handler: Home,
		},
		{
			Method:  "GET",
			Path:    "/sign-up",
			Handler: SignUp,
		},
		{
			Method:  "GET",
			Path:    "/log-in",
			Handler: LogIn,
		},
		{
			Method:  "GET",
			Path:    "/swagger/*",
			Handler: echoSwagger.WrapHandler,
		},
	}
}
