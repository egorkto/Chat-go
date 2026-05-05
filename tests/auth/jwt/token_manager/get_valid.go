package tests_auth_jwt_token_manager

import (
	"path/filepath"
	"time"

	auth_jwt_token_manager "github.com/egorkto/Chat-go/internal/auth/jwt/token_manager"
	"github.com/egorkto/Chat-go/internal/domain"
	tests_utils "github.com/egorkto/Chat-go/tests/utils"
)

func getTokenMangerValidConfig() auth_jwt_token_manager.Config {
	root := tests_utils.GetProjectRoot()

	privatePath := filepath.Join(root, "tests", "certs", "test_app.rsa")
	publicPath := filepath.Join(root, "tests", "certs", "test_app.rsa.pub")

	return auth_jwt_token_manager.Config{
		PrivatePath:    privatePath,
		PublicPath:     publicPath,
		AccessExpired:  15 * time.Minute,
		RefreshExpired: 24 * 7 * time.Hour,
		Issuer:         "test-issuer",
		Audience:       "test-audience",
	}
}

func getValidUser() domain.User {
	return domain.NewUser(
		1,
		1,
		"test-login",
		"Test",
	)
}
