package tests_auth_jwt_service

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	auth_jwt_service "github.com/egorkto/Chat-go/internal/auth/jwt/service"
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/egorkto/Chat-go/internal/validator"
	tests_mocks "github.com/egorkto/Chat-go/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSignUpValidation(t *testing.T) {
	var (
		mS_anythingArgs = []interface{}{mock.Anything, mock.Anything, mock.Anything}
		mG_anythingArgs = []interface{}{mock.Anything}
	)

	testCases := []struct {
		name      string
		full_name string
		login     string
		password  string
		wantErr   bool
	}{
		{
			name:      "Valid creds",
			full_name: gofakeit.Name(),
			login:     gofakeit.Username(),
			password:  gofakeit.Password(true, true, true, true, true, 10),
			wantErr:   false,
		},
		{
			name:      "Short full name",
			full_name: gofakeit.LetterN(2),
			login:     gofakeit.Username(),
			password:  gofakeit.Password(true, true, true, true, true, 10),
			wantErr:   true,
		},
		{
			name:      "Long full name",
			full_name: gofakeit.LetterN(110),
			login:     gofakeit.Username(),
			password:  gofakeit.Password(true, true, true, true, true, 10),
			wantErr:   true,
		},
		{
			name:      "Empty full name",
			full_name: "",
			login:     gofakeit.Username(),
			password:  gofakeit.Password(true, true, true, true, true, 10),
			wantErr:   true,
		},
		{
			name:      "Long password",
			full_name: gofakeit.Name(),
			login:     gofakeit.Username(),
			password:  gofakeit.Password(true, true, true, true, true, 120),
			wantErr:   true,
		},
		{
			name:      "Weak password",
			full_name: gofakeit.Name(),
			login:     gofakeit.Username(),
			password:  gofakeit.Password(true, false, true, false, false, 5),
			wantErr:   true,
		},
		{
			name:      "Empty password",
			full_name: gofakeit.Name(),
			login:     gofakeit.Username(),
			password:  "",
			wantErr:   true,
		},
		{
			name:      "Short login",
			full_name: gofakeit.Name(),
			login:     gofakeit.LetterN(2),
			password:  gofakeit.Password(true, true, true, true, true, 10),
			wantErr:   true,
		},
		{
			name:      "Long login",
			full_name: gofakeit.Name(),
			login:     gofakeit.LetterN(30),
			password:  gofakeit.Password(true, true, true, true, true, 10),
			wantErr:   true,
		},
		{
			name:      "Empty login",
			full_name: gofakeit.Name(),
			login:     "",
			password:  gofakeit.Password(true, true, true, true, true, 10),
			wantErr:   true,
		},
		{
			name:      "Empty creds",
			full_name: "",
			login:     "",
			password:  "",
			wantErr:   true,
		},
	}

	mockStorage := tests_mocks.NewUsersStorageMock()
	tokenManagerMock := tests_mocks.NewTokenManagerMock()

	mockStorage.On("CreateUser", mS_anythingArgs...).Return(domain.User{}, nil)
	tokenManagerMock.On("Generate", mG_anythingArgs...).Return(domain.JWT{}, nil)

	validator := validator.New()

	service := auth_jwt_service.New(tokenManagerMock, mockStorage, validator)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := domain.NewUser(0, 0, tc.full_name, tc.login)

			_, _, err := service.SignUp(context.Background(), user, tc.password)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				assert.NoError(t, err)
				mockStorage.AssertCalled(t, "CreateUser", mS_anythingArgs...)
				tokenManagerMock.AssertCalled(t, "Generate", mG_anythingArgs...)
			}
		})
	}
}
