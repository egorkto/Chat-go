package db

import (
	"fmt"
	"time"

	core_config "github.com/egorkto/Chat-go/internal/config"
)

type Config struct {
	Host     string        `envconfig:"HOST" required:"true"`
	Port     string        `envconfig:"PORT" default:"5432"`
	User     string        `envconfig:"USER" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Database string        `envconfig:"DB" required:"true"`
	Timeout  time.Duration `envconfig:"TIMEOUT" required:"true"`
}

func NewConfigMust() Config {
	config, err := core_config.New[Config]("POSTGRES")
	if err != nil {
		err = fmt.Errorf("get postgres config: %w", err)
		panic(err)
	}
	return config
}
