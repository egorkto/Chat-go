package chat_transport_websocket_hub

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	chat_transport "github.com/egorkto/Chat-go/internal/chat/transport"
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/gorilla/websocket"
)

type Hub struct {
	upgrader     websocket.Upgrader
	clients      map[Client]struct{}
	clientsMutex *sync.RWMutex
	broadcast    chan *chat_transport.MessageDTO
	saver        MessagesSaver
	saveQueue    chan *domain.Message
	log          *slog.Logger
}

type MessagesSaver interface {
	SaveMessage(ctx context.Context, msg domain.Message) error
}

func New(cfg Config, saver MessagesSaver, log *slog.Logger) *Hub {
	h := &Hub{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  cfg.ReadBufferSize,
			WriteBufferSize: cfg.WriteBufferSize,
		},
		clients:      make(map[Client]struct{}),
		clientsMutex: &sync.RWMutex{},
		broadcast:    make(chan *chat_transport.MessageDTO),
		saver:        saver,
		saveQueue:    make(chan *domain.Message, cfg.SaveBufferSize),
		log:          log,
	}

	go h.SaveWorker()

	return h
}

func (h *Hub) Upgrade(rw http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := h.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		return nil, fmt.Errorf("upgrade: %w", err)
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

			h.log.Debug("send broadcast message", slog.String("text", msg.Text))
		}

		h.clientsMutex.RUnlock()
	}
}

func (h *Hub) ReadFrom(conn *websocket.Conn, user domain.User) {
	client := Client{
		Conn: conn,
		User: user,
	}

	h.clientsMutex.Lock()
	h.clients[client] = struct{}{}
	h.clientsMutex.Unlock()

	defer func() {
		h.clientsMutex.Lock()
		delete(h.clients, client)
		h.clientsMutex.Unlock()
	}()

	for {
		msg := new(MessageInput)
		if err := client.Conn.ReadJSON(msg); err != nil {
			h.log.Debug("read json", slog.String("err", err.Error()))
			client.Conn.WriteJSON(
				ErrorResponse{
					Message: "failed to send message",
					Error:   "bad request",
				},
			)
		}

		h.log.Debug("recieved message", slog.String("text", msg.Text))

		domainMsg := domain.NewUninitializedMessage(
			client.User,
			msg.Text,
			time.Now(),
		)

		dto := chat_transport.DtoFromDomain(domainMsg)

		h.broadcast <- &dto
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
