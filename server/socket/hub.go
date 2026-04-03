package socket

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync/atomic"

	"github.com/coder/websocket"
)

type LobbyObserver interface {
	NotifyLobby(eventType string, data any)
}

type Hub struct {
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	clients    map[*Client]struct{}
	online     atomic.Int64
	observer   LobbyObserver
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
		broadcast:  make(chan []byte, 256),
		clients:    make(map[*Client]struct{}),
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for client := range h.clients {
				close(client.send)
				client.conn.CloseNow()
				delete(h.clients, client)
			}
			h.online.Store(0)
			return
		case client := <-h.register:
			h.clients[client] = struct{}{}
			h.online.Store(int64(len(h.clients)))
			data := map[string]any{
				"player": client.name,
				"online": len(h.clients),
			}
			h.emitSystem("player_joined", data)
			if h.observer != nil {
				h.observer.NotifyLobby("player_joined", data)
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.online.Store(int64(len(h.clients)))
				data := map[string]any{
					"player": client.name,
					"online": len(h.clients),
				}
				h.emitSystem("player_left", data)
				if h.observer != nil {
					h.observer.NotifyLobby("player_left", data)
				}
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
					client.conn.CloseNow()
				}
			}
		}
	}
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		slog.Error("websocket upgrade failed", "error", err)
		return
	}

	client := NewClient(h, conn)
	select {
	case h.register <- client:
	default:
		slog.Warn("hub register queue full")
		_ = conn.Close(websocket.StatusTryAgainLater, "server busy")
		return
	}

	go client.writePump(r.Context())
	client.readPump(r.Context())
}

func (h *Hub) SetObserver(obs LobbyObserver) {
	h.observer = obs
}

func (h *Hub) OnlineCount() int {
	return int(h.online.Load())
}

func (h *Hub) Broadcast(event OutboundEvent) {
	b, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal broadcast event", "error", err)
		return
	}
	select {
	case h.broadcast <- b:
	default:
		slog.Warn("broadcast channel full, dropping message", "type", event.Type)
	}
}

func (h *Hub) emitSystem(eventType string, data any) {
	h.Broadcast(OutboundEvent{
		Type: eventType,
		Data: data,
	})
}
