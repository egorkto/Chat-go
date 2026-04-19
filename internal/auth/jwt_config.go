package auth

import (
	"fmt"
	"time"

	"github.com/egorkto/Chat-go/internal/config"
)

type JWTConfig struct {
	PrivatePath    string        `envconfig:"PRIVATE_PATH" required:"true"`
	PublicPath     string        `envconfig:"PUBLIC_PATH" required:"true"`
	AccessExpired  time.Duration `envconfig:"ACCESS_EXP" required:"true"`
	RefreshExpired time.Duration `envconfig:"REFRESH_EXP" required:"true"`
	Issuer         string        `envconfig:"ISS" required:"true"`
	Audience       string        `envconfig:"AUD" required:"true"`
}

func NewJWTConfigMust() JWTConfig {
	config, err := config.New[JWTConfig]("JWT")
	if err != nil {
		err = fmt.Errorf("get jwt config: %w", err)
		panic(err)
	}
	return config
}
