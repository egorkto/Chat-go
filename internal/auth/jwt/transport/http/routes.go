package auth_jwt_transport_http

import (
	"net/http"

	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
	"github.com/labstack/echo/v5"
)

func (h HTTPHandler) Routes(mc transport_http.MiddlewaresContainer) []echo.Route {
	return []echo.Route{
		{
			Method:  http.MethodPost,
			Path:    "/sign-up",
			Handler: h.SignUp,
		},
		{
			Method:  http.MethodPost,
			Path:    "/log-in",
			Handler: h.LogIn,
		},
		{
			Method:  http.MethodPost,
			Path:    "/refresh",
			Handler: h.Refresh,
		},
	}
}
