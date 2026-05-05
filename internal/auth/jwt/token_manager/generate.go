package auth_jwt_token_manager

import (
	"fmt"
	"time"

	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

func (tm TokenManager) Generate(u domain.User) (domain.JWT, error) {
	subject := fmt.Sprintf("%s:%v", u.Login(), u.ID())

	accessIssued := jwt.NewNumericDate(time.Now())
	accessExpired := jwt.NewNumericDate(accessIssued.Add(tm.cfg.AccessExpired))

	refreshIssued := jwt.NewNumericDate(time.Now())
	refreshExpired := jwt.NewNumericDate(refreshIssued.Add(tm.cfg.RefreshExpired))

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
		return domain.JWT{}, fmt.Errorf(
			"signing access token: %w",
			err,
		)
	}

	refreshString, err := refreshToken.SignedString(tm.privateKey)
	if err != nil {
		return domain.JWT{}, fmt.Errorf(
			"signing refresh token: %w",
			err,
		)
	}

	jwt := domain.JWT{
		Access:  accessString,
		Refresh: refreshString,
	}

	return jwt, nil
}
