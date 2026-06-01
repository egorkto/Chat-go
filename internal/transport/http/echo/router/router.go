package transport_http_echo_router

import (
	"fmt"
	"log/slog"

	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Router struct {
	middlewares transport_http.MiddlewaresContainer
}

func New(mc transport_http.MiddlewaresContainer) Router {
	return Router{
		middlewares: mc,
	}
}

type HTTPHandler interface {
	Routes(mc transport_http.MiddlewaresContainer) []echo.Route
}

func (r Router) NewRouter(log *slog.Logger, handlers []HTTPHandler) (*echo.Echo, error) {
	e := echo.New()

	e.Logger = log

	e.Use(getRequestLoggerMiddleware())
	e.Use(middleware.Recover())

	reqVal, err := NewRequestValidator()
	if err != nil {
		return nil, fmt.Errorf("new request validator: %w", err)
	}
	e.Validator = reqVal

	e.HTTPErrorHandler = HTTPErrorHandler

	routes := []echo.Route{}
	for _, handler := range handlers {
		routes = append(routes, handler.Routes(r.middlewares)...)
	}

	addMany(e, routes)

	return e, nil
}

func getRequestLoggerMiddleware() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogMethod:    true,
		LogStatus:    true,
		LogRequestID: true,
		LogLatency:   true,
		LogHost:      true,
		HandleError:  false,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			c.Logger().Info("request",
				slog.String("uri", v.URI),
				slog.String("method", v.Method),
				slog.Int("status", v.Status),
				slog.String("request_id", v.RequestID),
				slog.Duration("latency", v.Latency),
				slog.String("host", v.Host),
			)
			return nil
		},
	})
}
