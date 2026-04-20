package command

import (
	"encoding/json"
	"log/slog"

	"github.com/sidinei-silva/eras-do-brasil-game/server/engine"
)

// Constantes para não usar strings soltas pelo código
const (
	TypeSendChat = "send_chat"
	TypeMove     = "move"
	TypeAttack   = "attack"
)

// A casca genérica que chega da rede
// ClientMessage representa o JSON cru enviado pelo cliente.
// O json.RawMessage guarda o resto do JSON como bytes para decodificar depois.
type ClientMessage struct {
	Type    string          `json:"type"`    // Ex: "send_chat"
	Payload json.RawMessage `json:"payload"` // O resto dos dados do comando
}

// PlayerCommand junta quem enviou com o que foi enviado.
type PlayerCommand struct {
	PlayerID string
	Message  ClientMessage
}

// ==========================================
// PAYLOADS ESPECÍFICOS (Crie um para cada ação)
// ==========================================

// Payload exclusivo para o Chat
type ChatPayload struct {
	Message string `json:"message"`
	Channel string `json:"channel"` // Ex: "global", "guild", "local"
}

// Payload exclusivo para Mover
type MovePayload struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// ProcessPlayerCommands roda no fluxo síncrono do tick
func ProcessPlayerCommands(engine *engine.GameLoop, cmds []PlayerCommand) {
	for _, cmd := range cmds {

		// Aqui fica o seu switch com as lógicas do JSON cru
		switch cmd.Message.Type {

		case "send_chat":
			// Você tem a ID do jogador e o Payload.
			var payload struct {
				Message string `json:"message"`
			}

			if err := json.Unmarshal(cmd.Message.Payload, &payload); err != nil {
				slog.Warn("Failed to parse chat command payload", "error", err)
				return // JSON inválido, ignora
			}

			slog.Info("Tick %d: Jogador %s enviou chat: %s", "tick", engine.TickCount(), "player_id", cmd.PlayerID, "message", payload.Message)

		case "move":
			slog.Info("Tick %d: Jogador %s tentou andar", "tick", engine.TickCount(), "player_id", cmd.PlayerID)

		default:
			slog.Warn("Comando desconhecido: %s", "type", cmd.Message.Type)
		}
	}
}
