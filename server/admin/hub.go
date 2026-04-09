package admin

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/sidinei-silva/eras-do-brasil-game/server/engine"
	"github.com/sidinei-silva/eras-do-brasil-game/server/world"
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
	mu      sync.Mutex
	clients map[*websocket.Conn]context.CancelFunc

	events chan []byte

	gameLoop          *engine.GameLoop
	world             *world.World
	playerOnlineCount func() int
}

func NewHub(gameLoop *engine.GameLoop, w *world.World, playerOnline func() int) *Hub {
	return &Hub{
		clients:           make(map[*websocket.Conn]context.CancelFunc),
		events:            make(chan []byte, 512),
		gameLoop:          gameLoop,
		world:             w,
		playerOnlineCount: playerOnline,
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			h.mu.Lock()
			for conn, cancel := range h.clients {
				cancel()
				conn.CloseNow()
			}
			h.clients = make(map[*websocket.Conn]context.CancelFunc)
			h.mu.Unlock()
			return

		case msg := <-h.events:
			h.mu.Lock()
			for conn, cancel := range h.clients {
				ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
				err := conn.Write(ctx2, websocket.MessageText, msg)
				cancel2()
				if err != nil {
					slog.Warn("admin client write failed, removing", "error", err)
					cancel()
					conn.CloseNow()
					delete(h.clients, conn)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) Publish(snap world.Snapshot) {
	h.send("gameloop", "tick", map[string]any{
		"tick":      snap.Tick,
		"game_time": snap.GameTime,
		"period":    snap.Period,
	})
}

func (h *Hub) NotifyLobby(eventType string, data any) {
	h.send("lobby", eventType, data)
}

func (h *Hub) send(category, eventType string, data any) {
	ev := newEvent(category, eventType, data)
	b, err := json.Marshal(ev)
	if err != nil {
		slog.Error("admin: failed to marshal event", "error", err)
		return
	}
	select {
	case h.events <- b:
	default:
		slog.Warn("admin event channel full, dropping", "type", eventType)
	}
}

func (h *Hub) sendToConn(conn *websocket.Conn, category, eventType string, data any) {
	ev := newEvent(category, eventType, data)
	b, err := json.Marshal(ev)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = conn.Write(ctx, websocket.MessageText, b)
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		slog.Error("admin websocket upgrade failed", "error", err)
		return
	}

	connCtx, connCancel := context.WithCancel(r.Context())

	h.mu.Lock()
	h.clients[conn] = connCancel
	h.mu.Unlock()

	slog.Info("admin connected", "total", h.clientCount())

	h.sendToConn(conn, "system", "welcome", map[string]any{
		"message":   "Admin console connected",
		"game_loop": h.gameLoop.Status(),
		"world":     h.world.Snapshot(),
	})

	h.readPump(connCtx, conn)

	h.mu.Lock()
	delete(h.clients, conn)
	h.mu.Unlock()
	connCancel()
	conn.CloseNow()

	slog.Info("admin disconnected", "total", h.clientCount())
}

func (h *Hub) readPump(ctx context.Context, conn *websocket.Conn) {
	conn.SetReadLimit(4096)
	for {
		readCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		_, message, err := conn.Read(readCtx)
		cancel()
		if err != nil {
			return
		}

		var cmd InboundCommand
		if err := json.Unmarshal(message, &cmd); err != nil {
			h.sendToConn(conn, "command", "error", map[string]string{
				"error": "invalid JSON",
			})
			continue
		}

		h.handleCommand(conn, cmd)
	}
}

func (h *Hub) clientCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.clients)
}
