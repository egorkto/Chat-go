package auth

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator struct {
	cfg        JWTConfig
	privateKey *rsa.PrivateKey
}

func NewJWTGenerator(cfg JWTConfig) (JWTGenerator, error) {
	privateBytes, err := os.ReadFile(cfg.PrivatePath)
	if err != nil {
		return JWTGenerator{}, fmt.Errorf("failed to read private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		return JWTGenerator{}, fmt.Errorf("failed to parse private key: %w", err)
	}

	return JWTGenerator{
		cfg:        cfg,
		privateKey: privateKey,
	}, nil
}

func (g JWTGenerator) Generate(u domain.User) (domain.JWT, error) {
	subject := fmt.Sprintf("%s:%v", u.FullName(), u.ID())

	accessIssued := jwt.NewNumericDate(time.Now())
	accessExpired := jwt.NewNumericDate(accessIssued.Add(g.cfg.AccessExpired))

	refreshIssued := jwt.NewNumericDate(time.Now())
	refreshExpired := jwt.NewNumericDate(refreshIssued.Add(g.cfg.RefreshExpired))

	accessClaims := jwt.RegisteredClaims{
		Issuer:    g.cfg.Issuer,
		Subject:   subject,
		Audience:  jwt.ClaimStrings{g.cfg.Audience},
		ExpiresAt: accessExpired,
		IssuedAt:  accessIssued,
	}

	refreshClaims := jwt.RegisteredClaims{
		Issuer:    g.cfg.Issuer,
		Subject:   subject,
		Audience:  jwt.ClaimStrings{g.cfg.Audience},
		ExpiresAt: refreshExpired,
		IssuedAt:  refreshIssued,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)

	accessString, err := accessToken.SignedString(g.privateKey)
	if err != nil {
		return domain.JWT{}, fmt.Errorf(
			"signing access token: %w",
			err,
		)
	}

	refreshString, err := refreshToken.SignedString(g.privateKey)
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
