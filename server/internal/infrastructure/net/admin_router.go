package net

import (
	"encoding/json"
	"log/slog"

	domainCmd "github.com/sidinei-silva/eras-do-brasil-game/server/internal/domain/command"
	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/infrastructure/engine"
)

// Defini uma interface para que o Router possa enviar mensagens de volta
// sem criar dependência circular com o Hub.
type AdminNotifier interface {
	Send(category, eventType string, data any) // O seu método h.send do hub.go
}

type AdminRouter struct {
	gameQueue *engine.CommandQueue
	notifier  AdminNotifier
	// worldManager *world.Manager // Injetar aqui se precisar ler dados do mundo instantaneamente
}

// NewAdminRouter cria o roteador e injeta a fila do motor do jogo
func NewAdminRouter(gameQueue *engine.CommandQueue, notifier AdminNotifier) *AdminRouter {
	return &AdminRouter{gameQueue: gameQueue, notifier: notifier}
}

// Route é chamado pelo readPump do cliente Admin
func (r *AdminRouter) Route(cmd domainCmd.PlayerCommand) {
	// O cmd.Message.Type diz-nos qual foi a ação pedida pelo Admin
	switch cmd.Message.Type {

	// ==========================================
	// 1. COMANDOS DE AÇÃO (Mudam o Mundo -> Vão para a Fila do GameLoop)
	// Substitui qualquer comando futuro como /spawn, /settime, /kill
	// ==========================================
	case "admin_set_time", "admin_spawn_npc", "admin_heal_player":
		r.gameQueue.Enqueue(cmd)
		slog.Info("Comando de Modo Deus enfileirado", "type", cmd.Message.Type)

	// ==========================================
	// 2. COMANDOS DE CONSULTA / OOB (Lêem o mundo instantaneamente)
	// Substitui o antigo comando: /npc <id>
	// ==========================================
	case "admin_get_npc":
		var payload struct {
			ID string `json:"id"`
		}

		if err := json.Unmarshal(cmd.Message.Payload, &payload); err != nil {
			r.notifier.Send("command", "error", map[string]string{"error": "invalid payload"})
			return
		}

		// Faz busca no WorldManager
		// state, ok := r.worldManager.GetNPCState(payload.ID)
		// if !ok {
		// 	r.notifier.Send("command", "error", map[string]string{"error": "npc not found", "id": payload.ID})
		// 	return
		// }

		// r.notifier.Send("command", "npc_details", state)
		slog.Info("Admin solicitou dados do NPC", "id", payload.ID)

	default:
		r.notifier.Send("command", "error", map[string]string{
			"error":   "unknown command",
			"command": cmd.Message.Type,
			"hint":    "Verifique a documentação da API Admin",
		})

	}
}
