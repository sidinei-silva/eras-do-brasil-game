package npc

import (
	"log/slog"
	"time"

	"github.com/sidinei-silva/eras-do-brasil-game/server/world"
)

type Manager struct {
	npcs []*Npc
}

func NewManager() (*Manager, error) {
	npcs, err := LoadNpcsFromFile()
	if err != nil {
		slog.Error("Erro ao carregar NPCs", "err", err)
		return nil, err
	}

	for _, npc := range npcs {
		slog.Debug("🆕🆕 ["+npc.Id+"] schedule loaded", "npc", npc.Name, "role", npc.Role, "blocks", len(npc.Schedule))
	}

	return &Manager{npcs: npcs}, nil
}

// ProcessTick é o coração do comportamento dos NPCs. É chamado pelo
// loop principal a cada tick. Para cada NPC, executa o ciclo:
//  1. Aplica decay contínuo (fome, fadiga, solidão sobem)
//  2. Verifica se atividade discreta atual terminou; se sim, aplica
//     efeito e decide próxima atividade
//  3. Em atividades contínuas, re-avalia decisão
//  4. Garante invariantes (clamp)
func (m *Manager) ProcessTick(gameTime world.GameTime, tickDuration time.Duration) {
	tickHours := tickDuration.Hours()

	for _, npc := range m.npcs {
		// Passo 1: decay sempre roda
		npcsInZone := m.getNpcsInZone(npc.CurrentZone, npc.Id)
		hasCompany := len(npcsInZone) > 0

		slog.Debug("ℹ️ ℹ️ ["+npc.Id+"] needs before decay",
			"id", npc.Id,
			"npc", npc.Name,
			"activity", npc.CurrentActivity,
			"has_company", hasCompany,
			"zone", npc.CurrentZone,
			"hunger", int(npc.Needs.Hunger),
			"fatigue", int(npc.Needs.Fatigue),
			"loneliness", int(npc.Needs.Loneliness),
			"score_hunger", npc.CalculateScores(gameTime.Time.Hour())["Hunger"],
			"score_fatigue", npc.CalculateScores(gameTime.Time.Hour())["Fatigue"],
			"score_schedule", npc.CalculateScores(gameTime.Time.Hour())["Schedule"],
			"npcs_in_zone_id", func() []string {
				ids := make([]string, 0, len(npcsInZone))
				for _, n := range npcsInZone {
					ids = append(ids, n.Id)
				}
				return ids
			}(),
		)

		npc.ApplyDecay(tickHours, hasCompany)

		// Passo 2 e 3: lógica de transição de atividade
		if npc.IsDiscreteActivity() {
			if npc.IsActivityComplete(gameTime) {
				npc.ApplyActivityEffects()
				slog.Info("✅✅ ["+npc.Id+"] activity done",
					"id", npc.Id,
					"npc", npc.Name,
					"activity", npc.CurrentActivity,
					"hunger", int(npc.Needs.Hunger),
					"fatigue", int(npc.Needs.Fatigue),
					"loneliness", int(npc.Needs.Loneliness),
				)
				m.decideNextActivity(npc, gameTime)
			}
			// Se não terminou, segue na atividade
		} else {
			// Atividades contínuas podem ser revisitadas a cada tick
			m.decideNextActivity(npc, gameTime)
		}

		// Passo 4: invariantes
		npc.ClampNeeds()

		slog.Debug("ℹ️ ℹ️ ["+npc.Id+"] needs updated",
			"id", npc.Id,
			"npc", npc.Name,
			"activity", npc.CurrentActivity,
			"zone", npc.CurrentZone,
			"hunger", int(npc.Needs.Hunger),
			"fatigue", int(npc.Needs.Fatigue),
			"loneliness", int(npc.Needs.Loneliness),
		)
	}
}

// decideNextActivity é o método do Manager que escolhe qual atividade
// o NPC deve fazer agora. Por enquanto usa lógica simples baseada em
// thresholds de need. Será substituído por scoring com pesos na próxima
// rodada (1.3 completo).
func (m *Manager) decideNextActivity(npc *Npc, gameTime world.GameTime) {

	hour := gameTime.Time.Hour()
	previousActivity := npc.CurrentActivity
	previousLocation := npc.CurrentZone
	desiredActivity, desiredLocation, winner, winnerScore := m.computeDesiredActivity(npc, hour)

	if desiredActivity == npc.CurrentActivity && desiredLocation == npc.CurrentZone {
		return
	}

	npc.TransitionTo(desiredActivity, desiredLocation, gameTime)

	slog.Info("🚶‍♂️‍➡️🚶‍♂️‍➡️ ["+npc.Id+"] transitioned",
		"id", npc.Id,
		"npc", npc.Name,
		"hour", hour,
		"from", previousActivity,
		"to", desiredActivity,
		"from_zone", previousLocation,
		"to_zone", desiredLocation,
		"winner", winner,
		"score", int(winnerScore),
	)
}

func (m *Manager) computeDesiredActivity(npc *Npc, hour int) (Activity, string, string, float64) {
	scores := npc.CalculateScores(hour)
	activity, location := ActivityIdle, npc.CurrentZone
	winner, maxScore := npc.PickWinnerScore(scores, npc.CurrentActivity)
	const minActionableScore = 10.0

	if maxScore < minActionableScore {
		return ActivityIdle, npc.CurrentZone, "none", maxScore
	}

	switch winner {
	case "Hunger":
		activity, location = pickMealForHour(hour), npc.EatingLocation
	case "Fatigue":
		activity, location = ActivitySleeping, npc.HomeLocation
	case "Schedule":
		action, _ := npc.ActiveScheduleAt(hour)
		activity, location = action.Activity, action.Location
	}

	return activity, location, winner, maxScore
}

func (m *Manager) GetAllNpcs() []*Npc {
	npcs := make([]*Npc, 0, len(m.npcs))
	for _, npc := range m.npcs {
		npcs = append(npcs, npc)
	}

	return npcs
}

func (m *Manager) GetNpcById(id string) (*Npc, bool) {
	for _, npc := range m.npcs {
		if npc.Id == id {
			return npc, true
		}
	}
	return nil, false
}

func (m *Manager) getNpcsInZone(zone string, npcId string) []*Npc {
	npcs := make([]*Npc, 0)
	for _, npc := range m.npcs {
		if npc.CurrentZone == zone && npc.Id != npcId {
			npcs = append(npcs, npc)
		}
	}
	return npcs
}
