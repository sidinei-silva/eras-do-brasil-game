package socket

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync/atomic"

	"github.com/coder/websocket"
)

// Hub é o coordenador central das conexões websocket do lobby.
// Ele recebe entradas por canais e processa tudo em uma única goroutine (Run),
// evitando locks explícitos para o mapa de clients.
type Hub struct {
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	clients    map[*Client]struct{}
	// online guarda uma cópia atômica do total para leitura fora da goroutine Run.
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

// Run é o loop principal do Hub.
// Toda mutação em clients acontece aqui para manter consistência de estado.
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// Shutdown gracioso: fecha canais de envio e força fechamento das conexões.
			for client := range h.clients {
				close(client.send)
				client.conn.CloseNow()
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
			// Remoção idempotente: só remove se ainda estiver no mapa.
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
					client.conn.CloseNow()
				}
			}
		}
	}
}

// ServeWS faz o upgrade HTTP → WebSocket e cria uma sessão Client.
// Bloqueia até a conexão ser encerrada — o goroutine do HTTP handler é o readPump.
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		// Aceita conexões de qualquer origem para facilitar testes locais.
		InsecureSkipVerify: true,
	})
	if err != nil {
		slog.Error("websocket upgrade failed", "error", err)
		return
	}

	client := NewClient(h, conn)
	// Evita bloquear a thread HTTP em caso de saturação de fila.
	select {
	case h.register <- client:
	default:
		slog.Warn("hub register queue full")
		_ = conn.Close(websocket.StatusTryAgainLater, "server busy")
		return
	}

	// writePump roda em goroutine separada (envia mensagens de saída).
	go client.writePump(r.Context())
	// readPump bloqueia aqui — quando a conexão cai, ServeWS retorna.
	// r.Context() fica vivo enquanto ServeWS estiver rodando.
	client.readPump(r.Context())
}

// OnlineCount devolve o total de conectados de forma thread-safe.
func (h *Hub) OnlineCount() int {
	return int(h.online.Load())
}

// Broadcast serializa um evento e enfileira para envio coletivo no loop do Hub.
// Usa select/default para nunca bloquear — importante porque emitSystem chama Broadcast
// de dentro de hub.Run, que é o único consumidor do canal. Bloquear seria deadlock.
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

// emitSystem envia eventos internos do lobby (join/leave etc).
func (h *Hub) emitSystem(eventType string, data any) {
	h.Broadcast(OutboundEvent{
		Type: eventType,
		Data: data,
	})
}
