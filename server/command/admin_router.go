package command

import (
	"encoding/json"
	"log/slog"

	"github.com/sidinei-silva/eras-do-brasil-game/server/config"
	"github.com/sidinei-silva/eras-do-brasil-game/server/npc"
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
	npcManager  *npc.Manager      // Injetar quando precisar de consultas imediatas no mundo
}

// NewAdminRouter cria o roteador e injeta a fila do motor do jogo
func NewAdminRouter(gameQueue *CommandQueue, notifier AdminNotifier, snapManager *snapshot.Manager, npcManager *npc.Manager) *AdminRouter {
	return &AdminRouter{gameQueue: gameQueue, notifier: notifier, snapManager: snapManager, npcManager: npcManager}
}

// Route é chamado pelo readPump da sessão admin
func (r *AdminRouter) Route(cmd PlayerCommand) {
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
	case "admin_get_npc_full":
		r.handleGetNpcFull(cmd)
	case "admin_get_snapshot":
		r.handleGetSnapshot()
	case "admin_get_npc_scores":
		r.handleGetNpcScores(cmd)
	case "admin_get_pois_in_block":
		r.handleGetPoisInBlock(cmd)
	case "admin_get_npc_knowledge":
		r.handleGetNpcKnowledge(cmd)

	default:
		r.notifier.Send("command", "error", map[string]string{
			"error":   "unknown command",
			"command": cmd.Message.Type,
			"hint":    "Verifique a documentação da API Admin",
		})
	}
}

func (r *AdminRouter) handleGetNpcFull(cmd PlayerCommand) {
	var payload struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(cmd.Message.Payload, &payload); err != nil {
		r.notifier.Send("command", "error", map[string]string{"error": "invalid payload"})
		return
	}

	snap := r.snapManager.GetSnapshot()
	if snap == nil {
		r.notifier.Send("command", "error", map[string]string{"error": "snapshot not available"})
		return
	}

	npc, found := snap.GetNPCById(payload.ID)
	if !found {
		r.notifier.Send("command", "error", map[string]string{"error": "npc not found", "id": payload.ID})
		return
	}

	type scheduleActionDTO struct {
		Activity  string `json:"activity"`
		PoiId     string `json:"poiId"`
		StartHour int    `json:"startHour"`
		EndHour   int    `json:"endHour"`
	}
	type needDTO struct {
		Hunger     float64 `json:"hunger"`
		Fatigue    float64 `json:"fatigue"`
		Loneliness float64 `json:"loneliness"`
	}
	type needWeightDTO struct {
		Hunger     float64 `json:"hunger"`
		Fatigue    float64 `json:"fatigue"`
		Loneliness float64 `json:"loneliness"`
		Schedule   float64 `json:"schedule"`
	}
	type npcFullDTO struct {
		ID              string              `json:"id"`
		Name            string              `json:"name"`
		Role            string              `json:"role"`
		CurrentBlock    string              `json:"currentBlock"`
		CurrentPoi      string              `json:"currentPoi"`
		CurrentActivity string              `json:"currentActivity"`
		Description     string              `json:"description"`
		Backstory       string              `json:"backstory"`
		HomePoi         string              `json:"homePoi"`
		EatingPoi       string              `json:"eatingPoi"`
		Needs           needDTO             `json:"needs"`
		NeedsWeight     needWeightDTO       `json:"needsWeight"`
		Schedule        []scheduleActionDTO `json:"schedule"`
	}

	schedule := make([]scheduleActionDTO, 0, len(npc.Schedule))
	for _, action := range npc.Schedule {
		schedule = append(schedule, scheduleActionDTO{
			Activity:  string(action.Activity),
			PoiId:     action.PoiId,
			StartHour: action.StartHour,
			EndHour:   action.EndHour,
		})
	}

	r.notifier.Send("command", "admin_get_npc_full", npcFullDTO{
		ID:              npc.Id,
		Name:            npc.Name,
		Role:            string(npc.Role),
		CurrentBlock:    npc.CurrentBlock,
		CurrentPoi:      npc.CurrentPoi,
		CurrentActivity: string(npc.CurrentActivity),
		Description:     npc.Description,
		Backstory:       npc.Backstory,
		HomePoi:         npc.HomePoi,
		EatingPoi:       npc.EatingPoi,
		Needs: needDTO{
			Hunger:     npc.Needs.Hunger,
			Fatigue:    npc.Needs.Fatigue,
			Loneliness: npc.Needs.Loneliness,
		},
		NeedsWeight: needWeightDTO{
			Hunger:     npc.NeedsWeight.Hunger,
			Fatigue:    npc.NeedsWeight.Fatigue,
			Loneliness: npc.NeedsWeight.Loneliness,
			Schedule:   npc.NeedsWeight.Schedule,
		},
		Schedule: schedule,
	})

	if config.Log.CommandRouting {
		slog.Debug("admin consultou NPC completo", "id", payload.ID)
	}
}

func (r *AdminRouter) handleGetSnapshot() {
	snap := r.snapManager.GetSnapshot()
	if snap == nil {
		r.notifier.Send("command", "error", map[string]string{"error": "snapshot not available"})
		return
	}
	r.notifier.Send("command", "admin_get_snapshot", snap)
	if config.Log.CommandRouting {
		slog.Debug("admin consultou snapshot", "tick", snap.Tick, "gameTime", snap.GetGameTime())
	}
}

func (r *AdminRouter) handleGetNpcScores(cmd PlayerCommand) {
	var payload struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(cmd.Message.Payload, &payload); err != nil {
		r.notifier.Send("command", "error", map[string]string{"error": "invalid payload"})
		return
	}

	snap := r.snapManager.GetSnapshot()
	if snap == nil {
		r.notifier.Send("command", "error", map[string]string{"error": "snapshot not available"})
		return
	}

	scores := snap.GetNpcScores(payload.ID)
	if scores == nil {
		r.notifier.Send("command", "error", map[string]string{"error": "npc scores not found", "id": payload.ID})
		return
	}

	r.notifier.Send("command", "admin_get_npc_scores", scores)
	if config.Log.CommandRouting {
		slog.Debug("admin consultou scores do NPC", "id", payload.ID)
	}
}

func (r *AdminRouter) handleGetPoisInBlock(cmd PlayerCommand) {
	var payload struct {
		BlockID string `json:"block_id"`
	}
	if err := json.Unmarshal(cmd.Message.Payload, &payload); err != nil {
		r.notifier.Send("command", "error", map[string]string{"error": "invalid payload"})
		return
	}

	snap := r.snapManager.GetSnapshot()
	if snap == nil {
		r.notifier.Send("command", "error", map[string]string{"error": "snapshot not available"})
		return
	}

	block, found := snap.GetBlockById(payload.BlockID)
	if !found {
		r.notifier.Send("command", "error", map[string]string{"error": "block not found", "id": payload.BlockID})
		return
	}

	poiCounts := snap.GetPoiCountsInBlock(payload.BlockID)
	r.notifier.Send("command", "admin_get_pois_in_block", map[string]interface{}{
		"block_id":   block.Id,
		"block_name": block.Name,
		"pois":       poiCounts,
	})
	if config.Log.CommandRouting {
		slog.Debug("admin consultou POIs do bloco", "block_id", payload.BlockID)
	}
}

func (r *AdminRouter) handleGetNpcKnowledge(cmd PlayerCommand) {
	var payload struct {
		ID string `json:"id"`
	}

	knowledgeConfig := r.npcManager.GetKnowledgeConfig()

	if err := json.Unmarshal(cmd.Message.Payload, &payload); err != nil {
		r.notifier.Send("command", "error", map[string]string{"error": "invalid payload"})
		return
	}

	snap := r.snapManager.GetSnapshot()
	if snap == nil {
		r.notifier.Send("command", "error", map[string]string{"error": "snapshot not available"})
		return
	}

	npc, found := snap.GetNPCById(payload.ID)
	if !found {
		r.notifier.Send("command", "error", map[string]string{"error": "npc not found", "id": payload.ID})
		return
	}

	knowledge := npc.GetActiveKnowledge(snap.GameTime.Time, knowledgeConfig)

	r.notifier.Send("command", "admin_get_npc_knowledge", knowledge)
	if config.Log.CommandRouting {
		slog.Debug("admin consultou conhecimento do NPC", "id", payload.ID, "returned", knowledge)
	}

}
