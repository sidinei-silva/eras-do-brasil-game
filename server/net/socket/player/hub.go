package net

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/coder/websocket"
	"github.com/sidinei-silva/eras-do-brasil-game/server/command"
)

// Hub mantém o conjunto de clientes ativos e transmite mensagens para eles.
type Hub struct {
	register    chan *PlayerSession
	unregister  chan *PlayerSession
	broadcast   chan []byte
	clients     map[*PlayerSession]struct{}
	online      atomic.Int64
	mu          sync.Mutex
	clientsByID map[string]*PlayerSession
	cmdQueue    *command.CommandQueue
	router      *command.PlayerRouter
}

func NewHub(queue *command.CommandQueue) *Hub {
	return &Hub{
		clients:    make(map[*PlayerSession]struct{}),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *PlayerSession),
		unregister: make(chan *PlayerSession),
		cmdQueue:   queue,
	}
}

// Injeta o router após a criação
func (h *Hub) SetRouter(r *command.PlayerRouter) {
	h.router = r
}

func (h *Hub) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			for client := range h.clients {
				close(client.send)
				client.conn.Close(websocket.StatusNormalClosure, "servidor encerrando conexões")
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
			slog.Info("jogador conectado", "player", client.name, "id", client.id, "online", h.online.Load())

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
				slog.Info("jogador desconectado", "player", client.name, "id", client.id, "online", h.online.Load())
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

func (h *Hub) ServeWS(mux *http.ServeMux, ctx context.Context, wg *sync.WaitGroup) {
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})

		if err != nil {
			slog.Error("websocket upgrade failed", "error", err)
			return
		}

		wg.Add(2)
		defer wg.Done()

		// Instancia a sessão do jogador com queue e router já injetados.
		client := NewPlayerSession(h, conn, h.cmdQueue, h.router)

		select {
		case h.register <- client:
		default:
			slog.Warn("hub register queue full")
			_ = conn.Close(websocket.StatusTryAgainLater, "server busy")
			return
		}

		go client.writePump(ctx)
		// readPump roda na goroutine atual para manter o ciclo de vida da conexão.
		client.readPump(ctx, wg)
	})
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

func (h *Hub) SendToClient(clientID string, data []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if client, ok := h.clientsByID[clientID]; ok {
		select {
		case client.send <- data:
		default:
			slog.Warn("client send buffer full, dropping message", "client", clientID)
		}
	}
}
