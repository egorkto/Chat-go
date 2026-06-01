package auth_jwt_token_manager_test

import (
	"path/filepath"
	"time"

	auth_jwt_token_manager "github.com/egorkto/Chat-go/internal/auth/jwt/token_manager"
	"github.com/egorkto/Chat-go/internal/utils"
)

func getTokenMangerValidConfig() auth_jwt_token_manager.Config {
	root := utils.GetProjectRoot()

	privatePath := filepath.Join(root, "certs", "test_app.rsa")
	publicPath := filepath.Join(root, "certs", "test_app.rsa.pub")

	return auth_jwt_token_manager.Config{
		PrivatePath:       privatePath,
		PublicPath:        publicPath,
		AccessExpiration:  15 * time.Minute,
		RefreshExpiration: 24 * 7 * time.Hour,
		Issuer:            "test-issuer",
		Audience:          "test-audience",
	}
}
