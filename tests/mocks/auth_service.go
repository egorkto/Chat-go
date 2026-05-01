package mocks

import (
	"context"

	"github.com/egorkto/Chat-go/internal/auth"
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
	creds auth.Credentials,
) (domain.User, domain.JWT, error) {
	args := m.Called(ctx, user, creds)
	return args.Get(0).(domain.User), args.Get(1).(domain.JWT), args.Error(2)
}
