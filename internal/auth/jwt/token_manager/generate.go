package auth_jwt_token_manager

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (tm TokenManager) Generate(userID int, userLogin string) (JWTPair, error) {
	subject := fmt.Sprintf("%v:%v", userID, userLogin)

	accessIssued := jwt.NewNumericDate(time.Now())
	accessExpired := jwt.NewNumericDate(accessIssued.Add(tm.cfg.AccessExpiration))

	refreshIssued := jwt.NewNumericDate(time.Now())
	refreshExpired := jwt.NewNumericDate(refreshIssued.Add(tm.cfg.RefreshExpiration))

	accessClaims := jwt.RegisteredClaims{
		Issuer:    tm.cfg.Issuer,
		Subject:   subject,
		Audience:  jwt.ClaimStrings{tm.cfg.Audience},
		ExpiresAt: accessExpired,
		IssuedAt:  accessIssued,
	}

	refreshClaims := jwt.RegisteredClaims{
		Issuer:    tm.cfg.Issuer,
		Subject:   subject,
		Audience:  jwt.ClaimStrings{tm.cfg.Audience},
		ExpiresAt: refreshExpired,
		IssuedAt:  refreshIssued,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)

	accessString, err := accessToken.SignedString(tm.privateKey)
	if err != nil {
		return JWTPair{}, fmt.Errorf(
			"signed access string: %w",
			err,
		)
	}

	refreshString, err := refreshToken.SignedString(tm.privateKey)
	if err != nil {
		return JWTPair{}, fmt.Errorf(
			"signed refresh string: %w",
			err,
		)
	}

	pair := JWTPair{
		Access: Token{
			Signed:    accessString,
			UserID:    userID,
			UserLogin: userLogin,
			ExpiredAt: accessExpired.Time,
		},
		Refresh: Token{
			Signed:    refreshString,
			UserID:    userID,
			UserLogin: userLogin,
			ExpiredAt: refreshExpired.Time,
		},
	}

	return pair, nil
}
