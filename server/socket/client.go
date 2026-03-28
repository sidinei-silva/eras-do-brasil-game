package socket

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// writeWait define timeout maximo para uma escrita no socket.
	writeWait = 10 * time.Second
	// pongWait define quanto tempo aceitamos ficar sem receber pong.
	pongWait = 60 * time.Second
	// pingPeriod e menor que pongWait para manter conexao viva.
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

// Client modela uma sessao websocket ativa no lobby.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	// send e a fila de saida consumida apenas pelo writePump.
	send chan []byte
	name string
}

// NewClient cria uma sessao com nome temporario ate o jogador definir o proprio nome.
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
		name: fmt.Sprintf("guest-%d", time.Now().UnixNano()%100000),
	}
}

// readPump e responsavel por ler mensagens da conexao e traduzi-las para comandos.
// Regra importante: somente readPump le da conexao.
func (c *Client) readPump() {
	defer func() {
		// Tenta remover do Hub sem bloquear em caso de shutdown.
		select {
		case c.hub.unregister <- c:
		default:
		}
		_ = c.conn.Close()
	}()

	// Limites e deadlines para robustez de conexao.
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// Fechamento inesperado vira warning para facilitar diagnostico.
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Warn("unexpected websocket close", "error", err)
			}
			break
		}

		var in InboundEvent
		if err := json.Unmarshal(message, &in); err != nil {
			// Mensagem invalida nao derruba conexao; so responde erro.
			c.sendJSON(OutboundEvent{Type: "error", Data: "invalid_json"})
			continue
		}

		c.handleInbound(in)
	}
}

// writePump e responsavel por enviar mensagens para o socket e manter heartbeat ping/pong.
// Regra importante: somente writePump escreve na conexao.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Canal fechado indica que sessao foi encerrada pelo Hub.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			// Ping periodico evita conexao zumbi atras de NAT/proxy.
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleInbound implementa os comandos basicos do lobby.
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
		c.hub.Broadcast(OutboundEvent{Type: "player_renamed", Data: map[string]string{"from": old, "to": c.name}})
	case "chat":
		// Publica mensagem de chat para todos conectados.
		body := strings.TrimSpace(in.Body)
		if body == "" {
			return
		}
		c.hub.Broadcast(OutboundEvent{Type: "chat", Data: map[string]string{"from": c.name, "body": body}})
	default:
		c.sendJSON(OutboundEvent{Type: "error", Data: "unknown_event_type"})
	}
}

// sendJSON enfileira uma resposta para envio assinc via writePump.
func (c *Client) sendJSON(event OutboundEvent) {
	b, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal outbound event", "error", err)
		return
	}

	select {
	case c.send <- b:
	default:
	}
}
