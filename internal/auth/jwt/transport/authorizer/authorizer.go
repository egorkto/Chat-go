package auth_jwt_transport_authorizer

import auth_jwt_token_manager "github.com/egorkto/Chat-go/internal/auth/jwt/token_manager"

type Authorizer struct {
	verifier TokenVerifier
}

type TokenVerifier interface {
	Verify(token string) (auth_jwt_token_manager.Token, error)
}

func New(i TokenVerifier) Authorizer {
	return Authorizer{
		verifier: i,
	}
}
