package command

import (
	"encoding/json"
	"log/slog"
)

type PlayerRouter struct {
	gameQueue *CommandQueue
	// chatService *chat.Service // Quando você criar o serviço de chat, ele entra aqui
}

func NewPlayerRouter(gameQueue *CommandQueue) *PlayerRouter {
	return &PlayerRouter{
		gameQueue: gameQueue,
	}
}

// Route é chamado pelo readPump
func (r *PlayerRouter) Route(cmd PlayerCommand) {
	switch cmd.Message.Type {
	// ==========================================
	// 1. OOB COMMANDS (Assíncronos, fora do GameLoop)
	// ==========================================

	case TypeSendChat:
		// Faz o parse para a struct estrita criada no domínio
		var payload ChatPayload
		if err := json.Unmarshal(cmd.Message.Payload, &payload); err != nil {
			slog.Warn("Falha ao ler payload de chat", "error", err, "playerID", cmd.PlayerID)
			return
		}

	// ==========================================
	// 2. SIMULATION COMMANDS (Vão para a Fila do GameLoop)
	// ==========================================

	case TypeMove, TypeAttack:
		// Aqui não fazemos o parse do payload. Jogamos a casca inteira para a fila do GameLoop e ele que decide o que fazer.
		// Isso é importante para manter a lógica de jogo centralizada e evitar duplicação de código.
		r.gameQueue.Enqueue(cmd)
		slog.Info("Comando de simulação enfileirado", "type", cmd.Message.Type)
	default:
		slog.Warn("Comando desconhecido recebido", "type", cmd.Message.Type, "playerID", cmd.PlayerID)

	}
}
