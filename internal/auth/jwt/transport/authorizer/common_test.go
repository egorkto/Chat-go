package auth_jwt_transport_authorizer_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v5"
)

func getTestContext() *echo.Context {
	e := echo.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/test",
		strings.NewReader(""))

	rec := httptest.NewRecorder()

	return e.NewContext(req, rec)
}
