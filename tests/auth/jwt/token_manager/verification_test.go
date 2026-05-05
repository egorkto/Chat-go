package tests_auth_jwt_token_manager

import (
	"path/filepath"
	"testing"
	"time"

	auth_jwt_token_manager "github.com/egorkto/Chat-go/internal/auth/jwt/token_manager"
	tests_utils "github.com/egorkto/Chat-go/tests/utils"
	"github.com/stretchr/testify/assert"
)

func TestVerify_Exp(t *testing.T) {
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

			cfg.AccessExpired = tc.accessExpired
			cfg.RefreshExpired = tc.refreshExpired

			tm, err := auth_jwt_token_manager.New(cfg)
			if err != nil {
				t.Fatalf("failed to create token manager: %s", err.Error())
			}

			token, err := tm.Generate(getValidUser())
			if err != nil {
				t.Fatalf("failed to generate tokens: %s", err.Error())
			}

			_, err = tm.Verify(token.Access)
			if tc.accessWantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			_, err = tm.Verify(token.Refresh)
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

	root := tests_utils.GetProjectRoot()
	fakePath := filepath.Join(root, "tests", "certs", "test_fake_app.rsa")

	cfg.PrivatePath = fakePath

	tm, err := auth_jwt_token_manager.New(cfg)
	if err != nil {
		t.Fatalf("failed to create token manager: %s", err.Error())
	}

	domainToken, err := tm.Generate(getValidUser())
	if err != nil {
		t.Fatalf("failed to generate tokens: %s", err.Error())
	}

	_, err = tm.Verify(domainToken.Access)
	assert.Error(t, err)

	_, err = tm.Verify(domainToken.Refresh)
	assert.Error(t, err)
}
