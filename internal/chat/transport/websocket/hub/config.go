package chat_transport_websocket_hub

import (
	"fmt"

	"github.com/egorkto/Chat-go/internal/utils"
)

type Config struct {
	ReadBufferSize      int `envconfig:"READ_BUFFER_SIZE" default:"1024"`
	WriteBufferSize     int `envconfig:"WRITE_BUFFER_SIZE" default:"1024"`
	SaveBufferSize      int `envconfig:"SAVE_BUFFER_SIZE" default:"512"`
	CientSendBufferSize int `envconfig:"CLIENT_SEND_BUFFER_SIZE" default:"256"`
}

func NewConfigMust() Config {
	cfg, err := utils.NewEnvConfig[Config]("WEBSOCKET_HUB")

	if err != nil {
		err = fmt.Errorf("get websocket hub config: %w", err)
		panic(err)
	}

	return cfg
}
