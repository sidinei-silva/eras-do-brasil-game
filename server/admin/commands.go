package admin

import (
	"strings"

	"github.com/coder/websocket"
)

// InboundCommand é o formato de comandos enviados pelo admin client.
// O client envia: {"command": "/status"} ou {"command": "/players"}
type InboundCommand struct {
	Command string `json:"command"`
}

// handleCommand roteia um comando para o handler correto.
func (h *Hub) handleCommand(conn *websocket.Conn, cmd InboundCommand) {
	raw := strings.TrimSpace(cmd.Command)
	if raw == "" {
		return
	}

	// Remove o "/" inicial se presente.
	raw = strings.TrimPrefix(raw, "/")

	parts := strings.Fields(raw)
	name := strings.ToLower(parts[0])

	switch name {
	case "status":
		h.cmdStatus(conn)
	case "world":
		h.cmdWorld(conn)
	case "players":
		h.cmdPlayers(conn)
	case "help":
		h.cmdHelp(conn)
	default:
		h.sendToConn(conn, "command", "error", map[string]string{
			"error":   "unknown command",
			"command": name,
			"hint":    "type /help for available commands",
		})
	}
}

// cmdStatus retorna o status do game loop.
func (h *Hub) cmdStatus(conn *websocket.Conn) {
	h.sendToConn(conn, "command", "status", h.gameLoop.Status())
}

// cmdWorld retorna o snapshot do mundo.
func (h *Hub) cmdWorld(conn *websocket.Conn) {
	h.sendToConn(conn, "command", "world", h.world.Snapshot())
}

// cmdPlayers retorna a contagem de jogadores online.
// Vai evoluir para listar jogadores individuais quando Player existir.
func (h *Hub) cmdPlayers(conn *websocket.Conn) {
	h.sendToConn(conn, "command", "players", map[string]any{
		"online": h.playerOnlineCount(),
	})
}

// cmdHelp lista comandos disponíveis.
func (h *Hub) cmdHelp(conn *websocket.Conn) {
	h.sendToConn(conn, "command", "help", map[string]any{
		"commands": []map[string]string{
			{"name": "/status", "desc": "Status do game loop (tick count, uptime, duration)"},
			{"name": "/world", "desc": "Snapshot do mundo (game time, periodo)"},
			{"name": "/players", "desc": "Jogadores conectados"},
			{"name": "/help", "desc": "Lista de comandos disponiveis"},
		},
	})
}
