package tests_auth

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/egorkto/Chat-go/internal/auth"
	auth_service "github.com/egorkto/Chat-go/internal/auth/service"
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/egorkto/Chat-go/tests/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type setupFunc func(
	mS *mocks.UsersStorageMock,
	mV *mocks.ValidatorMock,
	mG *mocks.JWTGeneratorMock,
)

func TestSignUpFlow(t *testing.T) {
	var (
		mS_anythingArgs = []interface{}{mock.Anything, mock.Anything, mock.Anything}
		mV_anythingArgs = []interface{}{mock.Anything, mock.Anything}
		mG_anythingArgs = []interface{}{mock.Anything}
		validUser       = domain.NewUser(0, 0, gofakeit.Name())
		invalidUser     = domain.User{}
	)

	testCases := []struct {
		name    string
		setup   setupFunc
		wantErr bool
		user    domain.User
	}{
		{
			name: "No errors",
			setup: func(
				mS *mocks.UsersStorageMock,
				mV *mocks.ValidatorMock,
				mG *mocks.JWTGeneratorMock,
			) {
				mV.On("ValidatePassword", mV_anythingArgs...).Return(nil)
				mG.On("Generate", mG_anythingArgs...).Return(domain.JWT{}, nil)
				mS.On("CreateUser", mS_anythingArgs...).Return(domain.User{}, nil)
			},
			wantErr: false,
			user:    validUser,
		},
		{
			name: "Invalid user",
			setup: func(
				mS *mocks.UsersStorageMock,
				mV *mocks.ValidatorMock,
				mG *mocks.JWTGeneratorMock,
			) {
			},
			wantErr: true,
			user:    invalidUser,
		},
		{
			name: "Password validation error",
			setup: func(
				mS *mocks.UsersStorageMock,
				mV *mocks.ValidatorMock,
				mG *mocks.JWTGeneratorMock,
			) {
				mV.On("ValidatePassword", mV_anythingArgs...).Return(
					errors.New("invalid password"))
			},
			wantErr: true,
			user:    validUser,
		},
		{
			name: "Creating user error",
			setup: func(
				mS *mocks.UsersStorageMock,
				mV *mocks.ValidatorMock,
				mG *mocks.JWTGeneratorMock,
			) {
				mV.On("ValidatePassword", mV_anythingArgs...).Return(nil)
				mS.On("CreateUser", mS_anythingArgs...).Return(
					domain.User{},
					errors.New("failed to create user"))
			},
			wantErr: true,
			user:    validUser,
		},
		{
			name: "Token generation error",
			setup: func(
				mS *mocks.UsersStorageMock,
				mV *mocks.ValidatorMock,
				mG *mocks.JWTGeneratorMock,
			) {
				mV.On("ValidatePassword", mV_anythingArgs...).Return(nil)
				mS.On("CreateUser", mS_anythingArgs...).Return(domain.User{}, nil)
				mG.On("Generate", mG_anythingArgs...).Return(
					domain.JWT{},
					errors.New("failed to generate token"),
				)

			},
			wantErr: true,
			user:    validUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storageMock := mocks.NewUsersStorageMock()
			jwtGeneratorMock := mocks.NewJWTGeneratorMock()
			validatorMock := mocks.NewValidatorMock()

			service := auth_service.New(jwtGeneratorMock, storageMock, validatorMock)

			tc.setup(storageMock, validatorMock, jwtGeneratorMock)
			user := tc.user

			_, _, err := service.SignUp(context.Background(), user, auth.Credentials{})

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
