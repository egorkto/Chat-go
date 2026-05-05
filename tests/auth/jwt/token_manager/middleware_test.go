package tests_auth_jwt_token_manager

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	auth_jwt_token_manager "github.com/egorkto/Chat-go/internal/auth/jwt/token_manager"
	tests_utils "github.com/egorkto/Chat-go/tests/utils"
	"github.com/labstack/echo/v5"
)

func TestMiddleware(t *testing.T) {
	validUser := getValidUser()

	validCfg := getTokenMangerValidConfig()
	fakeKeyCfg := getTokenMangerValidConfig()

	root := tests_utils.GetProjectRoot()

	fakePath := filepath.Join(root, "tests", "certs", "test_fake_app.rsa")

	fakeKeyCfg.PrivatePath = fakePath

	tm, err := auth_jwt_token_manager.New(validCfg)
	if err != nil {
		t.Fatalf("failed to create token manager: %s", err.Error())
	}

	validToken, err := tm.Generate(validUser)
	if err != nil {
		t.Fatalf("failed to generate valid token: %s", err.Error())
	}

	fakeTm, err := auth_jwt_token_manager.New(fakeKeyCfg)
	if err != nil {
		t.Fatalf("failed to create fake token manager: %s", err.Error())
	}

	fakeToken, err := fakeTm.Generate(validUser)
	if err != nil {
		t.Fatalf("failed to generate fake token: %s", err.Error())
	}

	testCases := []struct {
		name         string
		header       string
		expectedCode int
	}{
		{
			name:         "Valid header",
			header:       fmt.Sprintf("Bearer %s", validToken.Access),
			expectedCode: http.StatusOK,
		},
		{
			name:         "Header without 'Bearer'",
			header:       validToken.Access,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Empty header",
			header:       "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Fake token in header",
			header:       fmt.Sprintf("Bearer %s", fakeToken.Access),
			expectedCode: http.StatusUnauthorized,
		},
	}

	mw := tm.EchoMiddleware()
	testHandler := func(c *echo.Context) error {
		return nil
	}

	handler := mw(testHandler)

	e := echo.New()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/test",
				strings.NewReader(""))

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Add("Authorization", tc.header)

			rec := httptest.NewRecorder()

			ctx := e.NewContext(req, rec)

			var code int

			err := handler(ctx)
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					code = he.Code
				} else {
					t.Errorf("Unexpected error: %s", err.Error())
				}
			} else {
				code = rec.Code
			}

			if code != tc.expectedCode {
				t.Errorf("expected status: %d, got %d", tc.expectedCode, code)
			}
		})
	}

}
