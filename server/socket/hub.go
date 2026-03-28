package socket

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// CheckOrigin controla se o browser pode abrir WS a partir de outra origem.
	// Na fase inicial, aceitamos tudo para facilitar testes locais.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Hub e o coordenador central das conexoes websocket do lobby.
// Ele recebe entradas por canais e processa tudo em uma unica goroutine (Run),
// evitando locks explicitos para o mapa de clients.
type Hub struct {
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	clients    map[*Client]struct{}
	// online guarda uma copia atomica do total para leitura fora da goroutine Run.
	online atomic.Int64
}

// NewHub inicializa filas e o conjunto de clientes conectados.
func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
		broadcast:  make(chan []byte, 256),
		clients:    make(map[*Client]struct{}),
	}
}

// Run e o loop principal do Hub.
// Toda mutacao em clients acontece aqui para manter consistencia de estado.
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// Shutdown gracioso: fecha canais de envio e conexoes abertas.
			for client := range h.clients {
				close(client.send)
				_ = client.conn.Close()
				delete(h.clients, client)
			}
			h.online.Store(0)
			return
		case client := <-h.register:
			// Novo jogador entra no lobby.
			h.clients[client] = struct{}{}
			h.online.Store(int64(len(h.clients)))
			h.emitSystem("player_joined", map[string]any{
				"player": client.name,
				"online": len(h.clients),
			})
		case client := <-h.unregister:
			// Remocao idempotente: so remove se ainda estiver no mapa.
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.online.Store(int64(len(h.clients)))
				h.emitSystem("player_left", map[string]any{
					"player": client.name,
					"online": len(h.clients),
				})
			}
		case message := <-h.broadcast:
			// Fan-out: replica a mensagem para todos os clientes conectados.
			// Se o buffer de um cliente lotar, removemos para proteger o Hub.
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
					_ = client.conn.Close()
				}
			}
		}
	}
}

// ServeWS faz o upgrade HTTP -> WebSocket e cria uma sessao Client.
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket upgrade failed", "error", err)
		return
	}

	client := NewClient(h, conn)
	// Evita bloquear a thread HTTP em caso de saturacao de fila.
	select {
	case h.register <- client:
	default:
		slog.Warn("hub register queue full")
		_ = conn.Close()
		return
	}

	// Cada conexao ganha duas goroutines: leitura e escrita.
	go client.writePump()
	go client.readPump()
}

// OnlineCount devolve o total de conectados de forma thread-safe.
func (h *Hub) OnlineCount() int {
	return int(h.online.Load())
}

// Broadcast serializa um evento e enfileira para envio coletivo no loop do Hub.
func (h *Hub) Broadcast(event OutboundEvent) {
	b, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal broadcast event", "error", err)
		return
	}
	h.broadcast <- b
}

// emitSystem envia eventos internos do lobby (join/leave etc).
func (h *Hub) emitSystem(eventType string, data any) {
	h.Broadcast(OutboundEvent{
		Type: eventType,
		Data: data,
	})
}
