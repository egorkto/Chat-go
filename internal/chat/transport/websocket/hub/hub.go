package chat_transport_websocket_hub

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/gorilla/websocket"
)

type Hub struct {
	clients      map[Client]struct{}
	clientsMutex *sync.RWMutex
	broadcast    chan *domain.Message
	upgrader     websocket.Upgrader
	saver        MessagesSaver
	saveQueue    chan *domain.Message
	log          *slog.Logger
}

type MessagesSaver interface {
	SaveMessage(ctx context.Context, msg domain.Message) error
}

func New(saver MessagesSaver, log *slog.Logger, cfg Config) *Hub {
	h := &Hub{
		clients:      make(map[Client]struct{}),
		clientsMutex: &sync.RWMutex{},
		broadcast:    make(chan *domain.Message),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  cfg.ReadBufferSize,
			WriteBufferSize: cfg.WriteBufferSize,
		},
		saver:     saver,
		saveQueue: make(chan *domain.Message, cfg.SaveBufferSize),
		log:       log,
	}

	go h.SaveWorker()

	return h
}

func (h *Hub) Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, fmt.Errorf("upgrading http: %w", err)
	}

	return conn, nil
}

func (h *Hub) StartBroadcasting() {
	for msg := range h.broadcast {
		h.clientsMutex.RLock()

		for client := range h.clients {
			if err := client.Conn.WriteJSON(msg); err != nil {
				h.log.Error(
					"writing json",
					slog.String("err", err.Error()),
				)
			}
		}

		h.clientsMutex.RUnlock()
	}
}

func (h *Hub) ReadFrom(client Client) error {
	h.clientsMutex.Lock()
	h.clients[client] = struct{}{}
	h.clientsMutex.Unlock()

	defer func() {
		h.clientsMutex.Lock()
		delete(h.clients, client)
		h.clientsMutex.Unlock()
	}()

	for {
		msg := new(MessageDTO)
		if err := client.Conn.ReadJSON(msg); err != nil {
			return fmt.Errorf(
				"reading json, %s: %w",
				err.Error(),
				domain.ErrInvalidArgument,
			)
		}

		domainMsg := domain.NewUninitializedMessage(
			client.User,
			msg.Text,
			time.Now(),
		)

		h.broadcast <- &domainMsg
		h.saveQueue <- &domainMsg
	}
}

func (h Hub) SaveWorker() {
	for msg := range h.saveQueue {
		if err := h.saver.SaveMessage(context.Background(), *msg); err != nil {
			h.log.Error("save message", slog.String("err", err.Error()))
		}
	}
}
