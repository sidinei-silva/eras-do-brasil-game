package command

import (
	"encoding/json"
	"log/slog"

	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/domain/command"
	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/domain/state"
)

type Manager struct{}

func NewManager() *Manager { return &Manager{} }

// ProcessCommands roda no fluxo síncrono do tick
func (m *Manager) ProcessCommands(gs *state.GameState, cmds []command.PlayerCommand) {
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

			slog.Info("Tick %d: Jogador %s enviou chat: %s", "tick", gs.TickCount, "player_id", cmd.PlayerID, "message", payload.Message)

		case "move":
			slog.Info("Tick %d: Jogador %s tentou andar", "tick", gs.TickCount, "player_id", cmd.PlayerID)

		default:
			slog.Warn("Comando desconhecido: %s", "type", cmd.Message.Type)
		}
	}
}
