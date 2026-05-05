package logger

import (
	"fmt"

	"github.com/egorkto/Chat-go/internal/utils"
)

type Config struct {
	Level  string `envconfig:"LEVEL" default:"DEBUG"`
	Folder string `envconfig:"FOLDER" required:"true"`
}

func NewConfigMust() Config {
	cfg, err := utils.NewEnvConfig[Config]("LOGGER")

	if err != nil {
		err = fmt.Errorf("get logger config: %w", err)
		panic(err)
	}

	return cfg
}
