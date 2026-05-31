package auth_jwt_token_manager

import "time"

type JWTPair struct {
	Access  Token
	Refresh Token
}

type Token struct {
	Signed    string
	UserID    int
	UserLogin string
	ExpiredAt time.Time
}
