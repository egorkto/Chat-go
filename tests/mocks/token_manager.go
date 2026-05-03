package tests_mocks

import (
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type TokenManagerMock struct {
	mock.Mock
}

func NewTokenManagerMock() *TokenManagerMock {
	return &TokenManagerMock{}
}

func (m *TokenManagerMock) Generate(u domain.User) (domain.JWT, error) {
	args := m.Called(u)
	return args.Get(0).(domain.JWT), args.Error(1)
}

func (m *TokenManagerMock) Verify(token string) (*jwt.Token, error) {
	args := m.Called(token)
	return args.Get(0).(*jwt.Token), args.Error(1)
}
