package http_server

import (
	"fmt"
	"time"

	"github.com/egorkto/Chat-go/internal/config"
)

type Config struct {
	Port            int           `envconfig:"PORT" default:"DEBUG"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"20s"`
}

func NewConfigMust() Config {
	cfg, err := config.New[Config]("HTTP_SERVER")

	if err != nil {
		err = fmt.Errorf("get server config: %w", err)
		panic(err)
	}

	return cfg
}
