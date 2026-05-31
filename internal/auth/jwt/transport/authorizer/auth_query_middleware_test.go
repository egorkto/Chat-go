package auth_jwt_transport_authorizer_test

import (
	"testing"

	auth_jwt_token_manager "github.com/egorkto/Chat-go/internal/auth/jwt/token_manager"
	auth_jwt_token_manager_mocks "github.com/egorkto/Chat-go/internal/auth/jwt/token_manager/mocks"
	auth_jwt_transport_authorizer "github.com/egorkto/Chat-go/internal/auth/jwt/transport/authorizer"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/mock"
)

func TestAuthQueryMiddleware_Valid(t *testing.T) {
	handler := withAuthQueryMiddleware()

	ctx := getTestContext()

	params := ctx.Request().URL.Query()
	params.Set("token", "test-token")
	ctx.Request().URL.RawQuery = params.Encode()

	err := handler(ctx)
	if err != nil {
		t.Errorf("unexpected error: %s",
			err.Error(),
		)
	}
}

func withAuthQueryMiddleware() echo.HandlerFunc {
	tmMock := auth_jwt_token_manager_mocks.NewTokenManagerMock()
	tmMock.On("Verify", mock.Anything).Return(auth_jwt_token_manager.Token{}, nil)

	a := auth_jwt_transport_authorizer.New(tmMock)

	mw := a.AuthorizeQueryMiddleware()

	handler := func(c *echo.Context) error {
		return nil
	}

	return mw(handler)
}
