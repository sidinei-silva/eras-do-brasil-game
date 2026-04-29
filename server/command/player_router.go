package command

import (
	"encoding/json"
	"log/slog"
)

type PlayerRouter struct {
	gameQueue *CommandQueue
	// chatService *chat.Service // Injetar quando o serviço de chat existir.
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
	// 1. COMANDOS OOB (assíncronos, fora do GameLoop)
	// ==========================================

	case TypeSendChat:
		// Faz o parse para a struct estrita criada no domínio
		var payload ChatPayload
		if err := json.Unmarshal(cmd.Message.Payload, &payload); err != nil {
			slog.Warn("Falha ao ler payload de chat", "error", err, "playerID", cmd.PlayerID)
			return
		}

	// ==========================================
	// 2. COMANDOS DE SIMULAÇÃO (vão para a fila do GameLoop)
	// ==========================================

	case TypeMove, TypeAttack:
		// Não faz parse aqui: a decisão fica centralizada no fluxo síncrono do GameLoop.
		r.gameQueue.Enqueue(cmd)
		slog.Info("Comando de simulação enfileirado", "type", cmd.Message.Type)
	default:
		slog.Warn("Comando desconhecido recebido", "type", cmd.Message.Type, "playerID", cmd.PlayerID)

	}
}
