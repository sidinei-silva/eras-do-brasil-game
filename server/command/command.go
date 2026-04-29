package command

import (
	"encoding/json"
	"log/slog"
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
// PAYLOADS ESPECÍFICOS (um por ação)
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
func ProcessPlayerCommands(tickCount int64, cmds []PlayerCommand) {
	for _, cmd := range cmds {

		// Switch central para processar o JSON cru recebido da fila.
		switch cmd.Message.Type {

		case "send_chat":
			// Aqui temos a ID do jogador e o payload já recebido no comando.
			var payload struct {
				Message string `json:"message"`
			}

			if err := json.Unmarshal(cmd.Message.Payload, &payload); err != nil {
				slog.Warn("Failed to parse chat command payload", "error", err)
				return // JSON inválido, ignora
			}

			slog.Info("Jogador enviou chat", "tick", tickCount, "player_id", cmd.PlayerID, "message", payload.Message)

		case "move":
			slog.Info("Jogador tentou andar", "tick", tickCount, "player_id", cmd.PlayerID)

		default:
			slog.Warn("Comando desconhecido", "type", cmd.Message.Type)
		}
	}
}
