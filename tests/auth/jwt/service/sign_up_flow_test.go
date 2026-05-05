package tests_auth_jwt_service

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	auth_jwt_service "github.com/egorkto/Chat-go/internal/auth/jwt/service"
	"github.com/egorkto/Chat-go/internal/domain"
	tests_mocks "github.com/egorkto/Chat-go/tests/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type setupFunc func(
	mS *tests_mocks.UsersStorageMock,
	mV *tests_mocks.ValidatorMock,
	mG *tests_mocks.TokenManagerMock,
)

func TestSignUp_Flow(t *testing.T) {
	var (
		mS_anythingArgs = []interface{}{mock.Anything, mock.Anything, mock.Anything}
		mV_anythingArgs = []interface{}{mock.Anything, mock.Anything}
		mG_anythingArgs = []interface{}{mock.Anything}
		validUser       = domain.NewUser(0, 0, gofakeit.Name(), gofakeit.Username())
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
				mS *tests_mocks.UsersStorageMock,
				mV *tests_mocks.ValidatorMock,
				mG *tests_mocks.TokenManagerMock,
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
				mS *tests_mocks.UsersStorageMock,
				mV *tests_mocks.ValidatorMock,
				mG *tests_mocks.TokenManagerMock,
			) {
			},
			wantErr: true,
			user:    invalidUser,
		},
		{
			name: "Password validation error",
			setup: func(
				mS *tests_mocks.UsersStorageMock,
				mV *tests_mocks.ValidatorMock,
				mG *tests_mocks.TokenManagerMock,
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
				mS *tests_mocks.UsersStorageMock,
				mV *tests_mocks.ValidatorMock,
				mG *tests_mocks.TokenManagerMock,
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
				mS *tests_mocks.UsersStorageMock,
				mV *tests_mocks.ValidatorMock,
				mG *tests_mocks.TokenManagerMock,
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
			storageMock := tests_mocks.NewUsersStorageMock()
			tokenManagerMock := tests_mocks.NewTokenManagerMock()
			validatorMock := tests_mocks.NewValidatorMock()

			service := auth_jwt_service.New(tokenManagerMock, storageMock, validatorMock)

			tc.setup(storageMock, validatorMock, tokenManagerMock)
			user := tc.user

			_, _, err := service.SignUp(context.Background(), user, "validPassword123")

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
