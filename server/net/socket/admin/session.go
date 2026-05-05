package net

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/sidinei-silva/eras-do-brasil-game/server/command"
)

const (
	maxMessageSize = 1024
)

type AdminSession struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	id     string
	router *command.AdminRouter // Roteador de comandos do admin
}

func NewAdminSession(hub *Hub, conn *websocket.Conn, router *command.AdminRouter, ctx context.Context) *AdminSession {
	return &AdminSession{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),                 // Buffer de segurança
		id:     "admin-" + time.Now().Format("150405"), // Gera um ID provisório
		router: router,
	}
}

func (c *AdminSession) readPump(ctx context.Context, wg *sync.WaitGroup) {
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

		// Se tiver o AdminRouter, é aqui que o admin envia comandos:
		if c.router != nil {
			// 1. Decodifica o JSON base
			var clientMsg command.ClientMessage

			if err := json.Unmarshal(messageBytes, &clientMsg); err != nil {
				slog.Warn("Failed to unmarshal message", "client", c.id, "error", err)
				continue // Ignora JSON malformado
			}

			// 2. Monta o comando com a ID do jogador
			cmd := command.PlayerCommand{
				PlayerID: c.id, // ID desta sessão admin
				Message:  clientMsg,
			}

			// 3. Delega para o router decidir o próximo passo
			c.router.Route(cmd)
		}
	}
}

func (c *AdminSession) writePump(ctx context.Context) {
	defer c.conn.CloseNow()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// O hub fechou o canal
				c.conn.Close(websocket.StatusNormalClosure, "")
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
		case <-ctx.Done():
			slog.Debug("contexto cancelado, encerrando writePump do admin", "id", c.id)
			return
		}
	}
}
