package auth_jwt_token_manager

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

func (tm TokenManager) Verify(token string) (Token, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf(
				"unexpected signing method, %v: %w",
				token.Header["alg"],
				domain.ErrUnauthorized,
			)
		}
		return tm.publicKey, nil
	})
	if err != nil {
		return Token{}, fmt.Errorf("parse token, %s: %w", err.Error(), domain.ErrUnauthorized)
	}

	if !t.Valid {
		return Token{}, fmt.Errorf("invalid token: %w", domain.ErrUnauthorized)
	}

	login, id, err := getLoginIDFromToken(t)
	if err != nil {
		return Token{}, fmt.Errorf("get login, id from token: %w", err)
	}

	expiredAt, err := t.Claims.GetExpirationTime()
	if err != nil {
		return Token{}, fmt.Errorf("get expiration time, %s: %w", err.Error(), domain.ErrUnauthorized)
	}

	return Token{
		Signed:    token,
		UserID:    id,
		UserLogin: login,
		ExpiredAt: expiredAt.Time,
	}, err
}

func getLoginIDFromToken(token *jwt.Token) (string, int, error) {
	sub, err := token.Claims.GetSubject()
	if err != nil {
		return "", 0, fmt.Errorf(
			"get subject, %s: %w",
			err.Error(),
			domain.ErrUnauthorized,
		)
	}

	splited := strings.Split(sub, ":")
	if len(splited) != 2 {
		return "", 0, fmt.Errorf(
			"wrong subject format: %w",
			domain.ErrUnauthorized,
		)
	}

	strID := splited[0]
	login := splited[1]

	id, err := strconv.Atoi(strID)
	if err != nil {
		return "", 0, fmt.Errorf(
			"id is not a number, %s: %w",
			err.Error(),
			domain.ErrUnauthorized,
		)
	}

	return login, id, nil
}
