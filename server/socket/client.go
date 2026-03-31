package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/coder/websocket"
)

const (
	// writeWait define timeout máximo para uma escrita no socket.
	writeWait = 10 * time.Second
	// pongWait define quanto tempo aceitamos ficar sem receber pong.
	pongWait = 60 * time.Second
	// pingPeriod é menor que pongWait para manter conexão viva.
	pingPeriod = (pongWait * 9) / 10
	// maxMessageSize protege contra payloads grandes demais no lobby.
	maxMessageSize = 1024
)

// InboundEvent representa mensagens recebidas do cliente web.
type InboundEvent struct {
	Type string `json:"type"`
	Body string `json:"body,omitempty"`
	Name string `json:"name,omitempty"`
}

// OutboundEvent representa mensagens enviadas pelo servidor ao cliente.
type OutboundEvent struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

// Client modela uma sessão websocket ativa no lobby.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	// send é a fila de saída consumida apenas pelo writePump.
	send chan []byte
	name string
}

// NewClient cria uma sessão com nome temporário até o jogador definir o próprio nome.
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
		name: fmt.Sprintf("guest-%d", time.Now().UnixNano()%100000),
	}
}

// readPump é responsável por ler mensagens da conexão e traduzi-las para comandos.
// Regra importante: somente readPump lê da conexão.
func (c *Client) readPump(ctx context.Context) {
	defer func() {
		// Tenta remover do Hub sem bloquear em caso de shutdown.
		select {
		case c.hub.unregister <- c:
		default:
		}
		c.conn.CloseNow()
	}()

	c.conn.SetReadLimit(maxMessageSize)

	for {
		// Renova o deadline a cada mensagem recebida — equivalente ao pong handler do gorilla.
		readCtx, cancel := context.WithTimeout(ctx, pongWait)
		_, message, err := c.conn.Read(readCtx)
		cancel()
		if err != nil {
			status := websocket.CloseStatus(err)
			if status != websocket.StatusNormalClosure && status != websocket.StatusGoingAway {
				slog.Warn("unexpected websocket close", "error", err)
			}
			break
		}

		var in InboundEvent
		if err := json.Unmarshal(message, &in); err != nil {
			// Mensagem inválida não derruba conexão; só responde erro.
			c.sendJSON(OutboundEvent{Type: "error", Data: "invalid_json"})
			continue
		}

		c.handleInbound(in)
	}
}

// writePump é responsável por enviar mensagens para o socket e manter heartbeat ping/pong.
// Regra importante: somente writePump escreve na conexão.
func (c *Client) writePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.CloseNow()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// Canal fechado pelo shutdown do Hub — não precisa fechar explicitamente,
				// o defer CloseNow() cuida disso.
				return
			}
			writeCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Write(writeCtx, websocket.MessageText, message)
			cancel()
			if err != nil {
				return
			}
		case <-ticker.C:
			// Ping periódico evita conexão zumbi atrás de NAT/proxy.
			pingCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Ping(pingCtx)
			cancel()
			if err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// handleInbound implementa os comandos básicos do lobby.
func (c *Client) handleInbound(in InboundEvent) {
	switch strings.ToLower(strings.TrimSpace(in.Type)) {
	case "set_name":
		// Atualiza nome do jogador e notifica o lobby.
		name := strings.TrimSpace(in.Name)
		if name == "" {
			c.sendJSON(OutboundEvent{Type: "error", Data: "name_required"})
			return
		}
		old := c.name
		c.name = name
		c.sendJSON(OutboundEvent{Type: "name_updated", Data: map[string]string{"name": c.name}})
		renameData := map[string]string{"from": old, "to": c.name}
		c.hub.Broadcast(OutboundEvent{Type: "player_renamed", Data: renameData})
		if c.hub.observer != nil {
			c.hub.observer.NotifyLobby("player_renamed", renameData)
		}
	case "chat":
		// Publica mensagem de chat para todos conectados.
		body := strings.TrimSpace(in.Body)
		if body == "" {
			return
		}
		chatData := map[string]string{"from": c.name, "body": body}
		c.hub.Broadcast(OutboundEvent{Type: "chat", Data: chatData})
		if c.hub.observer != nil {
			c.hub.observer.NotifyLobby("chat", chatData)
		}
	default:
		c.sendJSON(OutboundEvent{Type: "error", Data: "unknown_event_type"})
	}
}

// sendJSON enfileira uma resposta para envio assíncrono via writePump.
func (c *Client) sendJSON(event OutboundEvent) {
	b, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal outbound event", "error", err)
		return
	}

	select {
	case c.send <- b:
	default:
		// Buffer cheio = cliente muito lento para processar respostas.
		// Mensagem descartada por backpressure — o hub remove o cliente se persistir.
		slog.Warn("client send buffer full, dropping message", "client", c.name, "type", event.Type)
	}
}
