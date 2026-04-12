package net

import (
	"encoding/json"
	"log/slog"

	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/domain/command"
	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/infrastructure/engine"
)

type CommandRouter struct {
	gameQueue *engine.CommandQueue
	// chatService *chat.Service // Quando você criar o serviço de chat, ele entra aqui
}

func NewCommandRouter(gameQueue *engine.CommandQueue) *CommandRouter {
	return &CommandRouter{
		gameQueue: gameQueue,
	}
}

// Route é chamado pelo readPump
func (r *CommandRouter) Route(cmd command.PlayerCommand) {
	switch cmd.Message.Type {
	// ==========================================
	// 1. OOB COMMANDS (Assíncronos, fora do GameLoop)
	// ==========================================

	case command.TypeSendChat:
		// Faz o parse para a struct estrita criada no domínio
		var payload command.ChatPayload
		if err := json.Unmarshal(cmd.Message.Payload, &payload); err != nil {
			slog.Warn("Falha ao ler payload de chat", "error", err, "playerID", cmd.PlayerID)
			return
		}

	// ==========================================
	// 2. SIMULATION COMMANDS (Vão para a Fila do GameLoop)
	// ==========================================

	case command.TypeMove, command.TypeAttack:
		// Aqui não fazemos o parse do payload. Jogamos a casca inteira para a fila do GameLoop e ele que decide o que fazer.
		// Isso é importante para manter a lógica de jogo centralizada e evitar duplicação de código.
		r.gameQueue.Enqueue(cmd)
		slog.Info("Comando de simulação enfileirado", "type", cmd.Message.Type)
	default:
		slog.Warn("Comando desconhecido recebido", "type", cmd.Message.Type, "playerID", cmd.PlayerID)

	}
}
