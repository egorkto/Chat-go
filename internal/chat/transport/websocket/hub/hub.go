package chat_transport_websocket_hub

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	chat_transport "github.com/egorkto/Chat-go/internal/chat/transport"
	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	user domain.User
	send chan *chat_transport.MessageDTO
}

type Hub struct {
	upgrader  websocket.Upgrader
	broadcast chan *chat_transport.MessageDTO

	register   chan *Client
	unregister chan *Client

	clients              map[*Client]struct{}
	clientSaveBufferSize int

	saver     MessagesSaver
	saveQueue chan *domain.Message

	log *slog.Logger
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
		broadcast:            make(chan *chat_transport.MessageDTO),
		register:             make(chan *Client),
		unregister:           make(chan *Client),
		clients:              make(map[*Client]struct{}),
		clientSaveBufferSize: cfg.SaveBufferSize,
		saver:                saver,
		saveQueue:            make(chan *domain.Message, cfg.SaveBufferSize),
		log:                  log,
	}

	return h
}

func (h *Hub) Run(ctx context.Context) {
	go h.saveWorker(ctx)

	for {
		select {
		case <-ctx.Done():
			for c := range h.clients {
				close(c.send)
				delete(h.clients, c)
			}
			return

		case c := <-h.register:
			h.clients[c] = struct{}{}

		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				close(c.send)
				delete(h.clients, c)
			}

		case msg := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					close(c.send)
					delete(h.clients, c)
				}
			}
		}
	}
}

func (h *Hub) Upgrade(rw http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := h.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		return nil, fmt.Errorf("upgrade: %w", err)
	}

	return conn, nil
}

func (h *Hub) ReadFrom(conn *websocket.Conn, user domain.User) {
	c := &Client{
		conn: conn,
		user: user,
		send: make(chan *chat_transport.MessageDTO, h.clientSaveBufferSize),
	}

	h.register <- c
	defer func() { h.unregister <- c }()

	go h.writePump(c)

	for {
		msg := new(MessageInput)
		if err := c.conn.ReadJSON(msg); err != nil {
			h.log.Debug("read json", slog.String("err", err.Error()))
			return
		}

		h.log.Debug("received message", slog.String("text", msg.Text))

		domainMsg := domain.NewUninitializedMessage(c.user, msg.Text, time.Now())
		dto := chat_transport.DtoFromDomain(domainMsg)

		h.broadcast <- &dto
		h.saveQueue <- &domainMsg
	}
}

func (h *Hub) writePump(c *Client) {
	defer c.conn.Close()

	for msg := range c.send {
		if err := c.conn.WriteJSON(msg); err != nil {
			h.log.Error("write json", slog.String("err", err.Error()))
			return
		}

		h.log.Debug("sent message", slog.String("text", msg.Text))
	}
}

func (h *Hub) saveWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-h.saveQueue:
			if !ok {
				return
			}
			if err := h.saver.SaveMessage(ctx, *msg); err != nil {
				h.log.Error("save message", slog.String("err", err.Error()))
			}
		}
	}
}
