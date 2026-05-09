package command

import (
	"encoding/json"
	"log/slog"

	"github.com/sidinei-silva/eras-do-brasil-game/server/config"
	"github.com/sidinei-silva/eras-do-brasil-game/server/snapshot"
)

// Interface para o Router enviar respostas sem criar dependência circular com o Hub.
type AdminNotifier interface {
	Send(category, eventType string, data any) // Método implementado pelo Hub Admin
}

type AdminRouter struct {
	gameQueue   *CommandQueue
	notifier    AdminNotifier
	snapManager *snapshot.Manager // Injetar para acessar o snapshot mais recente e garantir consistência com o que o admin vê no cliente.
	// npcManager *npc.Manager // Injetar quando precisar de consultas imediatas no mundo
}

// NewAdminRouter cria o roteador e injeta a fila do motor do jogo
func NewAdminRouter(gameQueue *CommandQueue, notifier AdminNotifier, snapManager *snapshot.Manager) *AdminRouter {
	return &AdminRouter{gameQueue: gameQueue, notifier: notifier, snapManager: snapManager}
}

// Route é chamado pelo readPump da sessão admin
func (r *AdminRouter) Route(cmd PlayerCommand) {
	// cmd.Message.Type define a ação solicitada pelo admin.
	switch cmd.Message.Type {

	// ==========================================
	// 1. COMANDOS DE AÇÃO (Mudam o Mundo -> Vão para a Fila do GameLoop)
	// Substitui qualquer comando futuro como /spawn, /settime, /kill
	// ==========================================
	case "admin_set_time", "admin_spawn_npc", "admin_heal_player":
		r.gameQueue.Enqueue(cmd)
		slog.Info("Comando de Modo Deus enfileirado", "type", cmd.Message.Type)

	// ==========================================
	// 2. COMANDOS DE CONSULTA / OOB (leem o mundo instantaneamente)
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

		snapshot := r.snapManager.GetSnapshot()
		if snapshot == nil {
			r.notifier.Send("command", "error", map[string]string{"error": "snapshot not available"})
			return
		}

		npc, found := snapshot.GetNPCById(payload.ID)

		if !found {
			r.notifier.Send("command", "error", map[string]string{"error": "npc not found", "id": payload.ID})
			return
		}

		r.notifier.Send("command", "admin_get_npc", npc)
		if config.Log.CommandRouting {
			slog.Debug("admin consultou NPC", "id", payload.ID)
		}

	// COMANDO GET SNAP
	case "admin_get_snapshot":
		snapshot := r.snapManager.GetSnapshot()
		if snapshot == nil {
			r.notifier.Send("command", "error", map[string]string{"error": "snapshot not available"})
			return
		}
		r.notifier.Send("command", "admin_get_snapshot", snapshot)
		if config.Log.CommandRouting {
			slog.Debug("admin consultou snapshot", "tick", snapshot.Tick, "gameTime", snapshot.GetGameTime())
		}

	case "admin_get_npc_scores":
		var payload struct {
			ID string `json:"id"`
		}

		if err := json.Unmarshal(cmd.Message.Payload, &payload); err != nil {
			r.notifier.Send("command", "error", map[string]string{"error": "invalid payload"})
			return
		}

		snapshot := r.snapManager.GetSnapshot()
		if snapshot == nil {
			r.notifier.Send("command", "error", map[string]string{"error": "snapshot not available"})
			return
		}

		scores := snapshot.GetNpcScores(payload.ID)
		if scores == nil {
			r.notifier.Send("command", "error", map[string]string{"error": "npc scores not found", "id": payload.ID})
			return
		}

		r.notifier.Send("command", "admin_get_npc_scores", scores)
		if config.Log.CommandRouting {
			slog.Debug("admin consultou scores do NPC", "id", payload.ID)
		}

	case "admin_get_pois_in_block":
		var payload struct {
			BlockID string `json:"block_id"`
		}

		if err := json.Unmarshal(cmd.Message.Payload, &payload); err != nil {
			r.notifier.Send("command", "error", map[string]string{"error": "invalid payload"})
			return
		}

		snapshot := r.snapManager.GetSnapshot()
		if snapshot == nil {
			r.notifier.Send("command", "error", map[string]string{"error": "snapshot not available"})
			return
		}

		block, found := snapshot.GetBlockById(payload.BlockID)
		if !found {
			r.notifier.Send("command", "error", map[string]string{"error": "block not found", "id": payload.BlockID})
			return
		}

		poiCounts := snapshot.GetPoiCountsInBlock(payload.BlockID)
		r.notifier.Send("command", "admin_get_pois_in_block", map[string]interface{}{
			"block_id":   block.Id,
			"block_name": block.Name,
			"pois":       poiCounts,
		})
		if config.Log.CommandRouting {
			slog.Debug("admin consultou POIs do bloco", "block_id", payload.BlockID)
		}

	default:
		r.notifier.Send("command", "error", map[string]string{
			"error":   "unknown command",
			"command": cmd.Message.Type,
			"hint":    "Verifique a documentação da API Admin",
		})

	}
}
