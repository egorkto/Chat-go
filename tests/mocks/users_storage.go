package mocks

import (
	"context"

	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/stretchr/testify/mock"
)

type UsersStorageMock struct {
	mock.Mock
}

func NewUsersStorageMock() *UsersStorageMock {
	return &UsersStorageMock{}
}

func (m *UsersStorageMock) CreateUser(
	ctx context.Context,
	user domain.User,
	pass string,
) (domain.User, error) {
	args := m.Called(user, pass)
	return args.Get(0).(domain.User), args.Error(1)
}
