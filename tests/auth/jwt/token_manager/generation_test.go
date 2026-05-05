package tests_auth_jwt_token_manager

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	auth_jwt_token_manager "github.com/egorkto/Chat-go/internal/auth/jwt/token_manager"
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGenerate_Subject(t *testing.T) {
	cfg := getTokenMangerValidConfig()

	tm, err := auth_jwt_token_manager.New(cfg)
	if err != nil {
		t.Fatalf("failed to create token manager: %s", err.Error())
	}

	validUser := getValidUser()

	domainJWT, err := tm.Generate(validUser)
	if err != nil {
		t.Fatalf("generating JWT: %v", err)
	}

	accessToken := domainJWT.Access
	err = testToken(accessToken, validUser)
	assert.NoError(t, err, "failed to test access token")

	refreshToken := domainJWT.Refresh
	err = testToken(refreshToken, validUser)
	assert.NoError(t, err, "failed to test refresh token")
}

func testToken(token string, user domain.User) error {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return fmt.Errorf("invalid token format: %s", token)
	}

	payload, _ := base64.RawURLEncoding.DecodeString(parts[1])

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return fmt.Errorf(
			"failed to unmarshal payload %s: %w",
			payload,
			err,
		)
	}

	sub, ok := claims["sub"]
	if !ok {
		return fmt.Errorf("subject claim not found in token payload: %s", payload)
	}

	splitedSub := strings.Split(sub.(string), ":")
	if len(splitedSub) != 2 {
		return fmt.Errorf(
			"invalid subject format %s",
			sub,
		)
	}

	issLogin := splitedSub[0]
	issID, err := strconv.Atoi(splitedSub[1])
	if err != nil {
		return fmt.Errorf(
			"Issuer id is not a number %s: %w",
			splitedSub[1],
			err,
		)
	}

	if issLogin != user.Login() || issID != user.ID() {
		return fmt.Errorf(
			"invalid issuer claims: expected %s:%d, got %s:%d",
			user.Login(),
			user.ID(),
			issLogin,
			issID,
		)
	}

	expClaim, ok := claims["exp"]
	if !ok {
		return fmt.Errorf("expiration claim not found in token payload: %s", payload)
	}

	exp, ok := expClaim.(float64)
	if !ok {
		return fmt.Errorf("Expiration claim is not a number: %s", expClaim)
	}

	iatClaim, ok := claims["iat"]
	if !ok {
		return fmt.Errorf("issued at claim not found in token payload: %s", payload)
	}

	iat, ok := iatClaim.(float64)
	if !ok {
		return fmt.Errorf("Issued at claim is not a number: %s", iatClaim)
	}

	if exp <= iat {
		return fmt.Errorf(
			"invalid expiration time: exp %f is not greater than iat %f",
			exp,
			iat,
		)
	}

	return nil
}
