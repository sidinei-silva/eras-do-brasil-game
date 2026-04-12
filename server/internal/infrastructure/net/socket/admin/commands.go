// package admin

// import (
// 	"strings"

// 	"github.com/coder/websocket"
// )

// type InboundCommand struct {
// 	Command string `json:"command"`
// }

// func (h *Hub) handleCommand(conn *websocket.Conn, cmd InboundCommand) {
// 	raw := strings.TrimSpace(cmd.Command)
// 	if raw == "" {
// 		return
// 	}

// 	raw = strings.TrimPrefix(raw, "/")

// 	parts := strings.Fields(raw)
// 	name := strings.ToLower(parts[0])

// 	switch name {
// 	case "status":
// 		h.cmdStatus(conn)
// 	case "world":
// 		h.cmdWorld(conn)
// 	case "players":
// 		h.cmdPlayers(conn)
// 	case "npcs":
// 		h.cmdNPCs(conn)
// 	case "npc":
// 		if len(parts) < 2 {
// 			h.sendToConn(conn, "command", "error", map[string]string{
// 				"error": "usage: /npc <id>",
// 			})
// 			return
// 		}
// 		h.cmdNPC(conn, parts[1])
// 	case "help":
// 		h.cmdHelp(conn)
// 	default:
// 		h.sendToConn(conn, "command", "error", map[string]string{
// 			"error":   "unknown command",
// 			"command": name,
// 			"hint":    "type /help for available commands",
// 		})
// 	}
// }

// func (h *Hub) cmdStatus(conn *websocket.Conn) {
// 	h.sendToConn(conn, "command", "status", h.gameLoop.Status())
// }

// func (h *Hub) cmdWorld(conn *websocket.Conn) {
// 	h.sendToConn(conn, "command", "world", h.world.Snapshot())
// }

// func (h *Hub) cmdPlayers(conn *websocket.Conn) {
// 	h.sendToConn(conn, "command", "players", map[string]any{
// 		"online": h.playerOnlineCount(),
// 	})
// }

// func (h *Hub) cmdNPCs(conn *websocket.Conn) {
// 	h.sendToConn(conn, "command", "npcs", map[string]any{
// 		"npcs": h.world.NPCs.AllStates(),
// 	})
// }

// func (h *Hub) cmdNPC(conn *websocket.Conn, id string) {
// 	state, ok := h.world.NPCs.State(id)
// 	if !ok {
// 		h.sendToConn(conn, "command", "error", map[string]string{
// 			"error": "npc not found",
// 			"id":    id,
// 		})
// 		return
// 	}
// 	h.sendToConn(conn, "command", "npc", state)
// }

//	func (h *Hub) cmdHelp(conn *websocket.Conn) {
//		h.sendToConn(conn, "command", "help", map[string]any{
//			"commands": []map[string]string{
//				{"name": "/status", "desc": "Status do game loop (tick count, uptime, duration)"},
//				{"name": "/world", "desc": "Snapshot do mundo (game time, periodo)"},
//				{"name": "/players", "desc": "Jogadores conectados"},
//				{"name": "/npcs", "desc": "Lista todos os NPCs com localização, atividade e necessidades"},
//				{"name": "/npc <id>", "desc": "Inspeciona um NPC específico pelo ID"},
//				{"name": "/help", "desc": "Lista de comandos disponiveis"},
//			},
//		})
//	}
package socket
