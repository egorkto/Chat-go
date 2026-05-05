package tests_mocks

import "github.com/stretchr/testify/mock"

type ValidatorMock struct {
	mock.Mock
}

func NewValidatorMock() *ValidatorMock {
	return &ValidatorMock{}
}

func (m *ValidatorMock) ValidatePassword(pass string, name string) error {
	args := m.Called(pass, name)
	return args.Error(0)
}
