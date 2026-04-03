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
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

type InboundEvent struct {
	Type string `json:"type"`
	Body string `json:"body,omitempty"`
	Name string `json:"name,omitempty"`
}

type OutboundEvent struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	name string
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
		name: fmt.Sprintf("guest-%d", time.Now().UnixNano()%100000),
	}
}

func (c *Client) readPump(ctx context.Context) {
	defer func() {
		select {
		case c.hub.unregister <- c:
		default:
		}
		c.conn.CloseNow()
	}()

	c.conn.SetReadLimit(maxMessageSize)

	for {
		_, message, err := c.conn.Read(ctx)
		if err != nil {
			status := websocket.CloseStatus(err)
			if status != websocket.StatusNormalClosure && status != websocket.StatusGoingAway {
				slog.Warn("unexpected websocket close", "error", err)
			}
			break
		}

		var in InboundEvent
		if err := json.Unmarshal(message, &in); err != nil {
			c.sendJSON(OutboundEvent{Type: "error", Data: "invalid_json"})
			continue
		}

		c.handleInbound(in)
	}
}

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
				return
			}
			writeCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Write(writeCtx, websocket.MessageText, message)
			cancel()
			if err != nil {
				return
			}
		case <-ticker.C:
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

func (c *Client) handleInbound(in InboundEvent) {
	switch strings.ToLower(strings.TrimSpace(in.Type)) {
	case "set_name":
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

func (c *Client) sendJSON(event OutboundEvent) {
	b, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal outbound event", "error", err)
		return
	}

	select {
	case c.send <- b:
	default:
		slog.Warn("client send buffer full, dropping message", "client", c.name, "type", event.Type)
	}
}
