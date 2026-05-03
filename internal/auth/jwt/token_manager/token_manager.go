package auth_jwt_token_manager

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	cfg        Config
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func New(cfg Config) (TokenManager, error) {
	privateBytes, err := os.ReadFile(cfg.PrivatePath)
	if err != nil {
		return TokenManager{}, fmt.Errorf("failed to read private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		return TokenManager{}, fmt.Errorf("failed to parse private key: %w", err)
	}

	publicBytes, err := os.ReadFile(cfg.PublicPath)
	if err != nil {
		return TokenManager{}, fmt.Errorf("failed to read public key: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		return TokenManager{}, fmt.Errorf("failed to parse public key: %w", err)
	}

	return TokenManager{
		cfg:        cfg,
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}
