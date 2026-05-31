package auth_jwt_transport_authorizer_test

import (
	"errors"
	"testing"

	auth_jwt_token_manager "github.com/egorkto/Chat-go/internal/auth/jwt/token_manager"
	auth_jwt_token_manager_mocks "github.com/egorkto/Chat-go/internal/auth/jwt/token_manager/mocks"
	auth_jwt_transport_authorizer "github.com/egorkto/Chat-go/internal/auth/jwt/transport/authorizer"
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/mock"
)

func TestAuthHeaderMiddleware_Valid(t *testing.T) {
	handler := withAuthHeaderMiddleware()

	ctx := getTestContext()

	ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx.Request().Header.Add("Authorization", "Bearer test-token")

	err := handler(ctx)
	if err != nil {
		t.Errorf("unexpected error %s",
			err.Error(),
		)
	}
}

func TestAuthHeaderMiddleware_WithoutPrefix(t *testing.T) {
	handler := withAuthHeaderMiddleware()

	ctx := getTestContext()

	ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx.Request().Header.Add("Authorization", "test-token")

	err := handler(ctx)
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Errorf("unexpected error type, expected: %s, have %v",
			domain.ErrUnauthorized.Error(),
			err,
		)
	}
}

func TestAuthHeaderMiddleware_WithoutHeader(t *testing.T) {
	handler := withAuthHeaderMiddleware()

	ctx := getTestContext()
	err := handler(ctx)
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Errorf("unexpected error type, expected: %s, have %v",
			domain.ErrUnauthorized.Error(),
			err,
		)
	}
}

func withAuthHeaderMiddleware() echo.HandlerFunc {
	tmMock := auth_jwt_token_manager_mocks.NewTokenManagerMock()
	tmMock.On("Verify", mock.Anything).Return(auth_jwt_token_manager.Token{}, nil)

	a := auth_jwt_transport_authorizer.New(tmMock)

	mw := a.AuthorizeHeaderMiddleware()

	handler := func(c *echo.Context) error {
		return nil
	}

	return mw(handler)
}
