package net

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/sidinei-silva/eras-do-brasil-game/server/command"
	"github.com/sidinei-silva/eras-do-brasil-game/server/snapshot"
)

type Event struct {
	Category  string `json:"category"`
	Type      string `json:"type"`
	Data      any    `json:"data,omitempty"`
	Timestamp string `json:"ts"`
}

func newEvent(category, eventType string, data any) Event {
	return Event{
		Category:  category,
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

type Hub struct {
	clients      map[*AdminSession]bool
	register     chan *AdminSession
	unregister   chan *AdminSession
	events       chan []byte             // Fila A: mensagens já serializadas (texto/erros)
	SnapshotChan chan *snapshot.Snapshot // Fila B: snapshots vindos do motor do jogo
	router       *command.AdminRouter
}

func NewHub() *Hub {
	return &Hub{
		clients:      make(map[*AdminSession]bool),
		register:     make(chan *AdminSession),
		unregister:   make(chan *AdminSession),
		events:       make(chan []byte, 512),
		SnapshotChan: make(chan *snapshot.Snapshot, 10), // Buffer de 10 ticks para não bloquear o GameLoop
	}
}

// Injeta o router após a criação
func (h *Hub) SetRouter(r *command.AdminRouter) {
	h.router = r
}

// ==========================================
// MÉTODOS DE ENTRADA DE DADOS
// ==========================================

// Usado para erros, respostas a comandos, chat.
// Faz o Marshal na hora porque não é chamado pelo GameLoop.
func (h *Hub) Send(category, eventType string, data any) {
	event := newEvent(category, eventType, data)

	payload, err := json.Marshal(event)

	if err != nil {
		slog.Error("admin: failed to marshal event", "error", err)
		return
	}

	select {
	case h.events <- payload:
	default:
		slog.Warn("admin event channel full, dropping", "type", eventType)
	}
}

// Usado EXCLUSIVAMENTE pelo GameLoop a cada 1 segundo.
// NÃO FAZ MARSHAL AQUI. Apenas entrega o ponteiro no canal e libera o motor.
// O select garante que se a fila do Admin estiver
// cheia (por latência de rede), o pacote é descartado e o motor do jogo não trava.
func (h *Hub) Publish(snap *snapshot.Snapshot) {
	select {
	case h.SnapshotChan <- snap:
	default:
		// Se o admin hub atrasar, o snapshot é descartado para preservar o ritmo do jogo.
	}
}

// ==========================================
// MÉTODOS DE PROCESSAMENTO E SAÍDA
// ==========================================
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

			return

		case client := <-h.register:
			h.clients[client] = true
			slog.Info("admin conectado", "id", client.id, "total_admins", len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				slog.Info("admin desconectado", "id", client.id, "total_admins", len(h.clients))
			}

		// FILA A: Distribui bytes que já vieram prontos do h.send()
		case payload := <-h.events:
			h.broadcastPayload(payload)

			// FILA B: O Hub recebe a foto do GameLoop, faz o Marshal pesado (longe do motor), e envia.
		case snap := <-h.SnapshotChan:
			event := newEvent("system", "snapshot", snap)
			payload, err := json.Marshal(event)

			if err == nil {
				h.broadcastPayload(payload)
			}

		}
	}
}

func (h *Hub) broadcastPayload(payload []byte) {
	for client := range h.clients {
		select {
		case client.send <- payload:
		default:
			// Se o buffer do cliente lotou (lag extremo), ele é desconectado para salvar o Hub
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) ServeWS(mux *http.ServeMux, ctx context.Context, wg *sync.WaitGroup) {
	mux.HandleFunc("/ws/admin", func(w http.ResponseWriter, r *http.Request) {

		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})

		if err != nil {
			slog.Error("admin websocket upgrade failed", "error", err)
			return
		}

		wg.Add(2)
		defer wg.Done()

		client := NewAdminSession(h, conn, h.router, ctx)

		select {
		case h.register <- client:
		default:
			slog.Warn("admin hub register queue full, recusando conexão", "id", client.id)
			_ = conn.Close(websocket.StatusTryAgainLater, "server busy")
			return
		}

		// ==========================================
		// 3. A MENSAGEM DE WELCOME DIRETA
		// ==========================================
		welcomeEvent := newEvent("system", "welcome", map[string]any{
			"message": "Admin console connected. Aguardando a sincronização do próximo tick...",
			// Não solicita snapshot aqui: o GameLoop publica automaticamente no SnapshotChan.
		})

		if payload, err := json.Marshal(welcomeEvent); err == nil {
			// Envia direto para este cliente, sem broadcast para os demais.
			client.send <- payload
		}

		go client.writePump(ctx)
		// readPump roda na goroutine atual para manter o ciclo de vida da conexão.
		client.readPump(ctx, wg)

	})
}
