package auth_jwt_token_manager_mocks

import (
	auth_jwt_token_manager "github.com/egorkto/Chat-go/internal/auth/jwt/token_manager"
	"github.com/stretchr/testify/mock"
)

type TokenManagerMock struct {
	mock.Mock
}

func NewTokenManagerMock() *TokenManagerMock {
	return &TokenManagerMock{}
}

func (m *TokenManagerMock) Generate(
	userID int,
	userLogin string,
) (auth_jwt_token_manager.JWTPair, error) {
	args := m.Called(userID, userLogin)
	return args.Get(0).(auth_jwt_token_manager.JWTPair), args.Error(1)
}

func (m *TokenManagerMock) Verify(token string) (auth_jwt_token_manager.Token, error) {
	args := m.Called(token)
	return args.Get(0).(auth_jwt_token_manager.Token), args.Error(1)
}
