package utils

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

func NewEnvConfig[T any](prefix string) (T, error) {
	var cfg T

	if err := envconfig.Process(prefix, &cfg); err != nil {
		return cfg, fmt.Errorf("process envconfig: %w", err)
	}

	return cfg, nil
}
