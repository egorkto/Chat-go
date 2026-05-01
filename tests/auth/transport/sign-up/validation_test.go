package tests_auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	auth_transport "github.com/egorkto/Chat-go/internal/auth/transport/http"
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/egorkto/Chat-go/tests/mocks"
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
				`{"full_name": "%s", "password": "%s"}`,
				gofakeit.Name(),
				gofakeit.Password(true, true, true, true, true, 10),
			),
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Short full name",
			json: fmt.Sprintf(
				`{"full_name": "%s", "password": "%s"}`,
				gofakeit.LetterN(1),
				gofakeit.Password(true, true, true, true, true, 10),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Long full name",
			json: fmt.Sprintf(
				`{"full_name": "%s", "password": "%s"}`,
				gofakeit.LetterN(40),
				gofakeit.Password(true, true, true, true, true, 10),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty full name",
			json: fmt.Sprintf(
				`{"full_name": "%s", "password": "%s"}`,
				"",
				gofakeit.Password(true, true, true, true, true, 10),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Null full name",
			json: fmt.Sprintf(
				`{"full_name": null, "password": "%s"}`,
				gofakeit.Password(true, true, true, true, true, 10),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "No full name",
			json: fmt.Sprintf(
				`{"password": "%s"}`,
				gofakeit.Password(true, true, true, true, true, 10),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Long password",
			json: fmt.Sprintf(
				`{"full_name": "%s", "password": "%s"}`,
				gofakeit.Name(),
				gofakeit.Password(true, true, true, true, true, 130),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Short password",
			json: fmt.Sprintf(
				`{"full_name": "%s", "password": "%s"}`,
				gofakeit.Name(),
				gofakeit.LetterN(4),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty password",
			json: fmt.Sprintf(
				`{"full_name": "%s", "password": "%s"}`,
				gofakeit.Name(),
				"",
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Null password",
			json: fmt.Sprintf(
				`{"full_name": "%s", "password": null}`,
				gofakeit.Name(),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "No password",
			json: fmt.Sprintf(
				`{"full_name": "%s"}`,
				gofakeit.Name(),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty creds",
			json: fmt.Sprintf(
				`{"full_name": "%s", "password": "%s"}`,
				"",
				"",
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Null creds",
			json:           `{"full_name": null, "password": null}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "No creds",
			json:           "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	serviceMock := mocks.NewAuthServiceMock()
	serviceMock.On(
		"SignUp",
		mock.Anything,
		mock.Anything,
		mock.Anything).Return(domain.User{}, domain.JWT{}, nil)

	transport := auth_transport.New(serviceMock)

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
