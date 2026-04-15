package player

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/domain/command"
	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/infrastructure/engine"
	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/infrastructure/net"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

type InboundEvent struct {
	Type string `json:"type"`
	Body string `json:"body,omitempty"`
}

type OutboundEvent struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

type PlayerSession struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	name     string
	id       string
	cmdQueue *engine.CommandQueue
	router   *net.CommandRouter
}

func NewPlayerSession(hub *Hub, conn *websocket.Conn, cmdQueue *engine.CommandQueue, router *net.CommandRouter) *PlayerSession {
	return &PlayerSession{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		name:     fmt.Sprintf("guest-%d", time.Now().UnixNano()%100000),
		id:       generateID(),
		cmdQueue: cmdQueue,
		router:   router,
	}
}

func (c *PlayerSession) readPump(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		select {
		case c.hub.unregister <- c:
		case <-ctx.Done():
		default:
		}

		c.conn.CloseNow()
	}()

	c.conn.SetReadLimit(maxMessageSize)

	for {
		_, messageBytes, err := c.conn.Read(ctx)

		if err != nil {
			status := websocket.CloseStatus(err)
			if status != websocket.StatusNormalClosure && status != websocket.StatusGoingAway {
				slog.Warn("unexpected websocket close", "error", err)
			}
			break
		}

		if c.router != nil {

			// 1. Decodifica o JSON base
			var clientMsg command.ClientMessage
			if err := json.Unmarshal(messageBytes, &clientMsg); err != nil {
				slog.Warn("Failed to unmarshal message", "client", c.name, "error", err)
				continue // Ignora JSON malformado
			}

			// 2. Monta o comando com a ID do jogador
			cmd := command.PlayerCommand{
				PlayerID: c.id, // A ID que você gerou para este client
				Message:  clientMsg,
			}

			// 3. Delega para o CommandRouter decidir o que fazer
			c.router.Route(cmd)
		}
	}
}

// writePump envia mensagens do Hub para a conexão WebSocket.
func (c *PlayerSession) writePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {

		case message, ok := <-c.send:
			if !ok {
				// O Hub fechou o canal, então fechamos a conexão
				c.conn.Close(websocket.StatusNormalClosure, "conexão encerrada pelo servidor")
				return
			}

			// Define um timeout para a escrita da mensagem no socket
			writeCtx, cancel := context.WithTimeout(ctx, time.Second*5)
			err := c.conn.Write(writeCtx, websocket.MessageText, message)
			cancel()

			if err != nil {
				slog.Error("erro de escrita", "error", err)
				return
			}

		case <-ticker.C:
			pingCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Ping(pingCtx)
			cancel()
			if err != nil {
				slog.Error("erro de ping", "error", err)
				return
			}

		case <-ctx.Done():
			slog.Info("contexto cancelado, encerrando writePump")
			return
		}
	}

}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
