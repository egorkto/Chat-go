package logger

import (
	"fmt"

	core_config "github.com/egorkto/Chat-go/internal/config"
)

type Config struct {
	Level  string `envconfig:"LEVEL" default:"DEBUG"`
	Folder string `envconfig:"FOLDER" required:"true"`
}

func NewConfigMust() Config {
	cfg, err := core_config.New[Config]("LOGGER")

	if err != nil {
		err = fmt.Errorf("get logger config: %w", err)
		panic(err)
	}

	return cfg
}
