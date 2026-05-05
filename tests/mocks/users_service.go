package tests_mocks

import (
	"context"

	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type UsersServiceMock struct {
	mock.Mock
}

func NewUsersServiceMock() *UsersServiceMock {
	return &UsersServiceMock{}
}

func (m *UsersServiceMock) CreateUser(
	ctx context.Context,
	user domain.User,
	pass string,
) (domain.User, error) {
	args := m.Called(user, pass)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UsersServiceMock) GetUserByLogin(
	ctx context.Context,
	login string,
) (domain.User, string, error) {
	args := m.Called(login)
	return args.Get(0).(domain.User), args.Get(1).(string), args.Error(2)
}

func (m *UsersServiceMock) GetUserByID(
	ctx context.Context,
	id int,
) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UsersServiceMock) GetUser(
	ctx context.Context,
	id int,
	token *jwt.Token,
) (domain.User, error) {
	args := m.Called(id, token)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UsersServiceMock) Refresh(
	ctx context.Context,
	refreshToken string,
) (string, error) {
	args := m.Called(refreshToken)
	return args.Get(0).(string), args.Error(1)
}
