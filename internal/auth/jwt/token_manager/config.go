package auth_jwt_token_manager

import (
	"fmt"
	"time"

	"github.com/egorkto/Chat-go/internal/utils"
)

type Config struct {
	PrivatePath    string        `envconfig:"PRIVATE_PATH" required:"true"`
	PublicPath     string        `envconfig:"PUBLIC_PATH" required:"true"`
	AccessExpired  time.Duration `envconfig:"ACCESS_EXP" required:"true"`
	RefreshExpired time.Duration `envconfig:"REFRESH_EXP" required:"true"`
	Issuer         string        `envconfig:"ISS" required:"true"`
	Audience       string        `envconfig:"AUD" required:"true"`
}

func NewConfigMust() Config {
	config, err := utils.NewEnvConfig[Config]("JWT")
	if err != nil {
		err = fmt.Errorf("get jwt config: %w", err)
		panic(err)
	}
	return config
}
