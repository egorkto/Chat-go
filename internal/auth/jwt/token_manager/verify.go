package auth_jwt_token_manager

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func (tm TokenManager) Verify(token string) (*jwt.Token, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tm.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !t.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return t, err
}
