package tests_auth_jwt

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	auth_jwt_transport_http "github.com/egorkto/Chat-go/internal/auth/jwt/transport/http"
	"github.com/egorkto/Chat-go/internal/domain"
	tests_mocks "github.com/egorkto/Chat-go/tests/mocks"
	"github.com/labstack/echo/v5"
	"github.com/presbrey/pkg/echovalidator"
	"github.com/stretchr/testify/mock"
)

func TestLogIn_Validation(t *testing.T) {
	validUser := domain.NewUser(1, 1, gofakeit.Name(), gofakeit.Username())
	validPassword := gofakeit.Password(true, true, true, true, true, 10)

	testCases := []struct {
		name           string
		json           string
		expectedStatus int
	}{
		{
			name: "Valid input",
			json: fmt.Sprintf(
				`{"login": "%s", "password": "%s"}`,
				validUser.Login(),
				validPassword,
			),
			expectedStatus: http.StatusOK,
		},
		{
			name: "Long login",
			json: fmt.Sprintf(
				`{"login": "%s", "password": "%s"}`,
				gofakeit.LetterN(40),
				validPassword,
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Long password",
			json: fmt.Sprintf(
				`{"login": "%s", "password": "%s"}`,
				validUser.Login(),
				gofakeit.Password(true, true, true, true, true, 200),
			),
			expectedStatus: http.StatusBadRequest,
		},
	}

	serviceMock := tests_mocks.NewAuthServiceMock()
	serviceMock.On(
		"LogIn",
		mock.Anything,
		mock.Anything,
		mock.Anything).Return(domain.User{}, domain.JWT{}, nil)
	serviceMock.On(
		"GetTokenExpires",
		mock.Anything).Return(time.Now(), nil)

	transport := auth_jwt_transport_http.New(serviceMock)

	e := echo.New()
	e.Validator = echovalidator.New()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/log-in",
				strings.NewReader(tc.json))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			ctx := e.NewContext(req, rec)

			var code int

			err := transport.LogIn(ctx)
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					code = he.Code
				} else {
					t.Errorf("Unexpected error: %s", err.Error())
				}
			} else {
				code = rec.Code
			}

			if code != tc.expectedStatus {
				t.Errorf("expected status: %d, got %d", tc.expectedStatus, code)
			}
		})
	}
}

func TestLogIn_InvalidUser(t *testing.T) {
	serviceMock := tests_mocks.NewAuthServiceMock()
	serviceMock.On(
		"LogIn",
		mock.Anything,
		mock.Anything,
		mock.Anything).Return(domain.User{}, domain.JWT{}, domain.ErrUnauthorized)
	serviceMock.On(
		"GetTokenExpires",
		mock.Anything).Return(time.Now(), nil)

	transport := auth_jwt_transport_http.New(serviceMock)

	e := echo.New()
	e.Validator = echovalidator.New()

	logInJSON := fmt.Sprintf(
		`{"login": "%s","password": "%s"}`,
		gofakeit.Username(),
		gofakeit.Password(true, true, true, true, true, 10),
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/log-in",
		strings.NewReader(logInJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	var code int

	err := transport.LogIn(ctx)
	if err != nil {
		if he, ok := err.(*echo.HTTPError); ok {
			t.Error(he.Message)
			code = he.Code
		} else {
			t.Errorf("Unexpected error: %s", err.Error())
		}
	} else {
		code = rec.Code
	}

	if code != http.StatusUnauthorized {
		t.Errorf("expected status: %d, got %d", http.StatusUnauthorized, code)
	}
}
