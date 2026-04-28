package tests_auth

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/egorkto/Chat-go/internal/auth"
	auth_service "github.com/egorkto/Chat-go/internal/auth/service"
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/egorkto/Chat-go/tests/mocks"
	"github.com/egorkto/Chat-go/validator"
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
		password  string
		wantErr   bool
	}{
		{
			name:      "Valid creds",
			full_name: gofakeit.Name(),
			password:  gofakeit.Password(true, true, true, true, true, 10),
			wantErr:   false,
		},
		{
			name:      "Short full name",
			full_name: gofakeit.LetterN(2),
			password:  gofakeit.Password(true, true, true, true, true, 10),
			wantErr:   true,
		},
		{
			name:      "Long full name",
			full_name: gofakeit.LetterN(40),
			password:  gofakeit.Password(true, true, true, true, true, 10),
			wantErr:   true,
		},
		{
			name:      "Empty full name",
			full_name: "",
			password:  gofakeit.Password(true, true, true, true, true, 10),
			wantErr:   true,
		},
		{
			name:      "Long password",
			full_name: gofakeit.Name(),
			password:  gofakeit.Password(true, true, true, true, true, 120),
			wantErr:   true,
		},
		{
			name:      "Weak password",
			full_name: gofakeit.Name(),
			password:  gofakeit.Password(true, false, true, false, false, 5),
			wantErr:   true,
		},
		{
			name:      "Empty password",
			full_name: gofakeit.Name(),
			password:  "",
			wantErr:   true,
		},
		{
			name:      "Empty creds",
			full_name: "",
			password:  "",
			wantErr:   true,
		},
	}

	mockStorage := mocks.NewUsersStorageMock()
	mockGenerator := mocks.NewJWTGeneratorMock()

	mockStorage.On("CreateUser", mS_anythingArgs...).Return(domain.User{}, nil)
	mockGenerator.On("Generate", mG_anythingArgs...).Return(domain.JWT{}, nil)

	validator := validator.New()

	service := auth_service.New(mockGenerator, mockStorage, validator)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := domain.NewUser(0, 0, tc.full_name)
			creds := auth.NewCredentials(tc.password)

			_, _, err := service.SignUp(context.Background(), user, creds)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				assert.NoError(t, err)
				mockStorage.AssertCalled(t, "CreateUser", mS_anythingArgs...)
				mockGenerator.AssertCalled(t, "Generate", mG_anythingArgs...)
			}
		})
	}
}
