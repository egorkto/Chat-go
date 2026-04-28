package mocks

import (
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/stretchr/testify/mock"
)

type JWTGeneratorMock struct {
	mock.Mock
}

func NewJWTGeneratorMock() *JWTGeneratorMock {
	return &JWTGeneratorMock{}
}

func (m *JWTGeneratorMock) Generate(u domain.User) (domain.JWT, error) {
	args := m.Called(u)
	return args.Get(0).(domain.JWT), args.Error(1)
}
