package transport_http_server

import (
	"fmt"
	"time"

	"github.com/egorkto/Chat-go/internal/utils"
)

type Config struct {
	Port            int           `envconfig:"PORT" default:"DEBUG"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"20s"`
}

func NewConfigMust() Config {
	cfg, err := utils.NewEnvConfig[Config]("HTTP_SERVER")

	if err != nil {
		err = fmt.Errorf("get server config: %w", err)
		panic(err)
	}

	return cfg
}
