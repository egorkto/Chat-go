package auth_jwt_token_manager_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	auth_jwt_token_manager "github.com/egorkto/Chat-go/internal/auth/jwt/token_manager"
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/egorkto/Chat-go/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestVerify_Expiration(t *testing.T) {
	testCases := []struct {
		name           string
		accessExpired  time.Duration
		refreshExpired time.Duration
		accessWantErr  bool
		refreshWantErr bool
	}{
		{
			name:           "Valid exps",
			accessExpired:  15 * time.Minute,
			refreshExpired: 24 * 7 * time.Hour,
			accessWantErr:  false,
			refreshWantErr: false,
		},
		{
			name:           "Zero access exp",
			accessExpired:  0,
			refreshExpired: 24 * 7 * time.Hour,
			accessWantErr:  true,
			refreshWantErr: false,
		},
		{
			name:           "Negative access exp",
			accessExpired:  -1 * time.Minute,
			refreshExpired: 24 * 7 * time.Hour,
			accessWantErr:  true,
			refreshWantErr: false,
		},
		{
			name:           "Zero refresh exp",
			accessExpired:  15 * time.Minute,
			refreshExpired: 0,
			accessWantErr:  false,
			refreshWantErr: true,
		},
		{
			name:           "Negative refresh exp",
			accessExpired:  15 * time.Minute,
			refreshExpired: -1 * time.Hour,
			accessWantErr:  false,
			refreshWantErr: true,
		},
		{
			name:           "Zero exps",
			accessExpired:  0,
			refreshExpired: 0,
			accessWantErr:  true,
			refreshWantErr: true,
		},
		{
			name:           "Negative exps",
			accessExpired:  -1 * time.Minute,
			refreshExpired: -1 * time.Hour,
			accessWantErr:  true,
			refreshWantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := getTokenMangerValidConfig()

			cfg.AccessExpiration = tc.accessExpired
			cfg.RefreshExpiration = tc.refreshExpired

			tm, err := auth_jwt_token_manager.New(cfg)
			if err != nil {
				t.Fatalf("failed to create token manager: %s", err.Error())
			}

			validUser := domain.NewUser(
				1,
				1,
				gofakeit.Name(),
				gofakeit.Username(),
			)

			token, err := tm.Generate(validUser.ID(), validUser.Login())
			if err != nil {
				t.Fatalf("failed to generate tokens: %s", err.Error())
			}

			_, err = tm.Verify(token.Access.Signed)
			if tc.accessWantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			_, err = tm.Verify(token.Refresh.Signed)
			if tc.refreshWantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVerify_FakeSign(t *testing.T) {
	cfg := getTokenMangerValidConfig()

	root := utils.GetProjectRoot()
	fakePath := filepath.Join(root, "certs", "test_fake_app.rsa")

	cfg.PrivatePath = fakePath

	tm, err := auth_jwt_token_manager.New(cfg)
	if err != nil {
		t.Fatalf("failed to create token manager: %s", err.Error())
	}

	validUser := domain.NewUser(
		1,
		1,
		gofakeit.Name(),
		gofakeit.Username(),
	)

	domainToken, err := tm.Generate(validUser.ID(), validUser.Login())
	if err != nil {
		t.Fatalf("failed to generate tokens: %s", err.Error())
	}

	_, err = tm.Verify(domainToken.Access.Signed)
	assert.Error(t, err)

	_, err = tm.Verify(domainToken.Refresh.Signed)
	assert.Error(t, err)
}
