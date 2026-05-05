package tests_mocks

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

func (m *UsersStorageMock) GetUserByLogin(
	ctx context.Context,
	login string,
) (domain.User, string, error) {
	args := m.Called(login)
	return args.Get(0).(domain.User), args.Get(1).(string), args.Error(2)
}

func (m *UsersStorageMock) GetUserByID(
	ctx context.Context,
	id int,
) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}
