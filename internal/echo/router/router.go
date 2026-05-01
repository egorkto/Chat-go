package echo_router

import (
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/presbrey/pkg/echovalidator"
)

func NewRouter(log *slog.Logger) *echo.Echo {
	e := echo.New()
	e.Logger = log
	e.Validator = echovalidator.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	return e
}
