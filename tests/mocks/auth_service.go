package tests_mocks

import (
	"context"
	"time"

	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/stretchr/testify/mock"
)

type AuthServiceMock struct {
	mock.Mock
}

func NewAuthServiceMock() *AuthServiceMock {
	return &AuthServiceMock{}
}

func (m *AuthServiceMock) SignUp(
	ctx context.Context,
	user domain.User,
	password string,
) (domain.User, domain.JWT, error) {
	args := m.Called(ctx, user, password)
	return args.Get(0).(domain.User), args.Get(1).(domain.JWT), args.Error(2)
}

func (m *AuthServiceMock) LogIn(
	ctx context.Context,
	login string,
	password string,
) (domain.User, domain.JWT, error) {
	args := m.Called(ctx, login, password)
	return args.Get(0).(domain.User), args.Get(1).(domain.JWT), args.Error(2)
}

func (m *AuthServiceMock) Refresh(
	ctx context.Context,
	refreshToken string,
) (domain.JWT, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(domain.JWT), args.Error(1)
}

func (m *AuthServiceMock) GetTokenExpires(token string) (time.Time, error) {
	args := m.Called(token)
	return args.Get(0).(time.Time), args.Error(1)
}
