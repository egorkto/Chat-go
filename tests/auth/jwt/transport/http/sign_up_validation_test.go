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

func TestValidation(t *testing.T) {
	testCases := []struct {
		name           string
		json           string
		expectedStatus int
	}{
		{
			name: "Valid input",
			json: fmt.Sprintf(
				`{"full_name": "%s", "login": "%s", "password": "%s"}`,
				gofakeit.Name(),
				gofakeit.Username(),
				gofakeit.Password(true, true, true, true, true, 10),
			),
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Short full name",
			json: fmt.Sprintf(
				`{"full_name": "%s", "login": "%s", "password": "%s"}`,
				gofakeit.LetterN(1),
				gofakeit.Username(),
				gofakeit.Password(true, true, true, true, true, 10),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Long full name",
			json: fmt.Sprintf(
				`{"full_name": "%s", "login": "%s", "password": "%s"}`,
				gofakeit.LetterN(120),
				gofakeit.Username(),
				gofakeit.Password(true, true, true, true, true, 10),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty full name",
			json: fmt.Sprintf(
				`{"full_name": "%s", "login": "%s", "password": "%s"}`,
				"",
				gofakeit.Username(),
				gofakeit.Password(true, true, true, true, true, 10),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Null full name",
			json: fmt.Sprintf(
				`{"full_name": null, "login": "%s", "password": "%s"}`,
				gofakeit.Username(),
				gofakeit.Password(true, true, true, true, true, 10),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "No full name",
			json: fmt.Sprintf(
				`{"login": "%s", "password": "%s"}`,
				gofakeit.Username(),
				gofakeit.Password(true, true, true, true, true, 10),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Long password",
			json: fmt.Sprintf(
				`{"full_name": "%s", "login": "%s", "password": "%s"}`,
				gofakeit.Name(),
				gofakeit.Username(),
				gofakeit.Password(true, true, true, true, true, 130),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Short password",
			json: fmt.Sprintf(
				`{"full_name": "%s", "login": "%s", "password": "%s"}`,
				gofakeit.Name(),
				gofakeit.Username(),
				gofakeit.LetterN(4),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty password",
			json: fmt.Sprintf(
				`{"full_name": "%s", "login": "%s", "password": "%s"}`,
				gofakeit.Name(),
				gofakeit.Username(),
				"",
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Null password",
			json: fmt.Sprintf(
				`{"full_name": "%s", "login": "%s", "password": null}`,
				gofakeit.Name(),
				gofakeit.Username(),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "No password",
			json: fmt.Sprintf(
				`{"full_name": "%s", "login": "%s", }`,
				gofakeit.Name(),
				gofakeit.Username(),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Long login",
			json: fmt.Sprintf(
				`{"full_name": "%s", "login": "%s", "password": "%s"}`,
				gofakeit.Name(),
				gofakeit.LetterN(40),
				gofakeit.Password(true, true, true, true, true, 10),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Short login",
			json: fmt.Sprintf(
				`{"full_name": "%s", "login": "%s", "password": "%s"}`,
				gofakeit.Name(),
				gofakeit.LetterN(2),
				gofakeit.Password(true, true, true, true, true, 10),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty login",
			json: fmt.Sprintf(
				`{"full_name": "%s", "login": "", "password": "%s"}`,
				gofakeit.Name(),
				gofakeit.Password(true, true, true, true, true, 10),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Null login",
			json: fmt.Sprintf(
				`{"full_name": "%s", "login": null, "password": "%s"}`,
				gofakeit.Name(),
				gofakeit.Password(true, true, true, true, true, 10),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "No login",
			json: fmt.Sprintf(
				`{"full_name": "%s", "password": "%s"}`,
				gofakeit.Name(),
				gofakeit.Password(true, true, true, true, true, 10),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty creds",
			json: fmt.Sprintf(
				`{"full_name": "%s", "login": "%s", "password": "%s"}`,
				"",
				"",
				"",
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Null creds",
			json:           `{"full_name": null, "login": null, "password": null}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "No creds",
			json:           "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	serviceMock := tests_mocks.NewAuthServiceMock()
	serviceMock.On(
		"SignUp",
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
				"/sign-up",
				strings.NewReader(tc.json))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			ctx := e.NewContext(req, rec)

			var code int

			err := transport.SignUp(ctx)
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
