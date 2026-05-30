package auth_jwt_token_manager_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
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

	validUser := domain.NewUser(
		1,
		1,
		gofakeit.Name(),
		gofakeit.Username(),
	)

	pair, err := tm.Generate(validUser.ID(), validUser.Login())
	if err != nil {
		t.Fatalf("generating JWT: %v", err)
	}

	err = testToken(pair.Access, validUser)
	assert.NoError(t, err, "failed to test access token")

	err = testToken(pair.Refresh, validUser)
	assert.NoError(t, err, "failed to test refresh token")
}

func testToken(token auth_jwt_token_manager.Token, user domain.User) error {
	if token.UserID != user.ID() {
		return fmt.Errorf("invalid token user id, expected: %s have: %s",
			user.ID(),
			token.UserID,
		)
	}

	if token.UserLogin != user.Login() {
		return fmt.Errorf("invalid token user login, expected: %s have: %s",
			user.Login(),
			token.UserLogin,
		)
	}

	tString := token.Signed
	parts := strings.Split(tString, ".")
	if len(parts) < 2 {
		return fmt.Errorf("invalid token format: %s", tString)
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

	subStr, ok := sub.(string)
	if !ok {
		return fmt.Errorf("subject claim is not a string: %s", sub)
	}

	expectedSub := fmt.Sprintf("%s:%d", user.Login(), user.ID())

	if subStr != expectedSub {
		return fmt.Errorf(
			"invalid subject claims: expected %s, got %s",
			expectedSub,
			subStr,
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

	expTime := time.Unix(int64(exp), 0)
	if expTime != token.ExpiredAt {
		return fmt.Errorf("signed token expired is not equal to token.ExpiredAt, %s:%s",
			expTime,
			token.ExpiredAt,
		)
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
