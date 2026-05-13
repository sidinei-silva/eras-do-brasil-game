package npc

import (
	"log/slog"
	"time"

	"github.com/sidinei-silva/eras-do-brasil-game/server/config"
	"github.com/sidinei-silva/eras-do-brasil-game/server/world"
)

type Manager struct {
	npcs            []*Npc
	knowledgeConfig KnowledgeConfig
}

func NewManager() (*Manager, error) {
	knowledgeConfig, err := LoadKnowledgeConfig()

	if err != nil {
		slog.Error("Erro ao carregar configuração de conhecimento", "err", err)
		return nil, err
	}

	npcs, err := LoadNpcsFromFile()
	if err != nil {
		slog.Error("Erro ao carregar NPCs", "err", err)
		return nil, err
	}

	if config.Log.NPCSchedule {
		for _, npc := range npcs {
			slog.Debug("npc schedule carregado", "id", npc.Id, "npc", npc.Name, "role", npc.Role, "blocos", len(npc.Schedule))
		}
	}

	return &Manager{npcs: npcs, knowledgeConfig: knowledgeConfig}, nil
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
	npcsInZoneByNPC := make(map[string][]*Npc, len(m.npcs))

	// Fase 1: atualizar estado individual e auto-gerar avistamentos.
	for _, npc := range m.npcs {
		npcsInZone := m.getNpcsInLocation(npc.CurrentBlock, npc.CurrentPoi, npc.Id)
		npcsInZoneByNPC[npc.Id] = npcsInZone

		npc.removeExpiredKnowledge(gameTime.Time, m.knowledgeConfig.ExpirationHours)

		for _, otherNpc := range npcsInZone {
			knowledge := NewKnowledgeAvistamentoNPC(otherNpc.Id, otherNpc.CurrentBlock, otherNpc.CurrentPoi, "direct", gameTime.Time)
			npc.addOrUpdateKnowledge(knowledge, gameTime.Time)
		}

		hasCompany := len(npcsInZone) > 0

		if config.Log.NPCNeeds {
			hour := gameTime.Time.Hour()
			scores := npc.CalculateScores(hour)
			ids := make([]string, 0, len(npcsInZone))
			for _, n := range npcsInZone {
				ids = append(ids, n.Id)
			}
			slog.Debug("npc needs antes do decay",
				"id", npc.Id,
				"npc", npc.Name,
				"activity", npc.CurrentActivity,
				"has_company", hasCompany,
				"zone", npc.CurrentBlock,
				"poi", npc.CurrentPoi,
				"hunger", int(npc.Needs.Hunger),
				"fatigue", int(npc.Needs.Fatigue),
				"loneliness", int(npc.Needs.Loneliness),
				"score_hunger", scores["Hunger"],
				"score_fatigue", scores["Fatigue"],
				"score_schedule", scores["Schedule"],
				"npcs_na_zona", ids,
			)
		}

		npc.ApplyDecay(tickHours, hasCompany)

		// Passo 2 e 3: lógica de transição de atividade
		if npc.IsDiscreteActivity() {
			if npc.IsActivityComplete(gameTime) {
				npc.ApplyActivityEffects()
				if config.Log.NPCBehavior {
					slog.Info("npc atividade concluída",
						"id", npc.Id,
						"npc", npc.Name,
						"activity", npc.CurrentActivity,
						"hunger", int(npc.Needs.Hunger),
						"fatigue", int(npc.Needs.Fatigue),
						"loneliness", int(npc.Needs.Loneliness),
					)
				}
				m.decideNextActivity(npc, gameTime)
			}
			// Se não terminou, segue na atividade
		} else {
			// Atividades contínuas podem ser revisitadas a cada tick
			m.decideNextActivity(npc, gameTime)
		}

		// Passo 4: invariantes
		npc.ClampNeeds()

		if config.Log.NPCNeeds {
			slog.Debug("npc needs após tick",
				"id", npc.Id,
				"npc", npc.Name,
				"activity", npc.CurrentActivity,
				"zone", npc.CurrentBlock,
				"hunger", int(npc.Needs.Hunger),
				"fatigue", int(npc.Needs.Fatigue),
				"loneliness", int(npc.Needs.Loneliness),
			)
		}
	}

	// Fase 2: fofoca. Usa o cache da primeira passada para não recalcular vizinhança.
	processedPairs := make(map[string]struct{})
	for _, npc := range m.npcs {
		npcsInZone := npcsInZoneByNPC[npc.Id]
		for _, otherNpc := range npcsInZone {
			pairKey := getPairKey(npc.Id, otherNpc.Id)
			if _, exists := processedPairs[pairKey]; exists {
				continue
			}
			processedPairs[pairKey] = struct{}{}

			lastGossipTime, hasGossiped := npc.LastGossipedWith[otherNpc.Id]
			if hasGossiped {
				cooldownDuration := time.Duration(m.knowledgeConfig.GossipCooldownHours) * time.Hour
				if gameTime.Time.Sub(lastGossipTime) < cooldownDuration {
					continue
				}
			}

			if !npc.exchangeKnowledgeWith(otherNpc, gameTime.Time, m.knowledgeConfig) {
				continue
			}

			if config.Log.NPCKnowledge {
				slog.Debug("Fofoca executada entre NPCs",
					"npc1_id", npc.Id,
					"npc1_name", npc.Name,
					"npc2_id", otherNpc.Id,
					"npc2_name", otherNpc.Name,
					"block", npc.CurrentBlock,
					"poi", npc.CurrentPoi,
				)
			}
		}
	}

}

// decideNextActivity é o método do Manager que escolhe qual atividade
// o NPC deve fazer agora. Por enquanto usa lógica simples baseada em
// thresholds de need. Será substituído por scoring com pesos na próxima
// rodada (1.3 completo).
func (m *Manager) decideNextActivity(npc *Npc, gameTime world.GameTime) {

	hour := gameTime.Time.Hour()
	previousActivity := npc.CurrentActivity
	previousLocation := npc.CurrentPoi

	desiredActivity, desiredLocation, winner, winnerScore := m.computeDesiredActivity(npc, hour)

	if desiredActivity == npc.CurrentActivity && desiredLocation == npc.CurrentPoi {
		return
	}

	npc.TransitionTo(desiredActivity, desiredLocation, gameTime)

	if config.Log.NPCBehavior {
		slog.Info("npc transicionou atividade",
			"id", npc.Id,
			"npc", npc.Name,
			"hora", hour,
			"de", previousActivity,
			"para", desiredActivity,
			"de_zona", previousLocation,
			"para_zona", desiredLocation,
			"motivo", winner,
			"score", int(winnerScore),
		)
	}
}

func (m *Manager) computeDesiredActivity(npc *Npc, hour int) (Activity, string, string, float64) {
	scores := npc.CalculateScores(hour)
	activity, location := ActivityIdle, npc.CurrentPoi
	winner, maxScore := npc.PickWinnerScore(scores, npc.CurrentActivity)
	const minActionableScore = 10.0

	if maxScore < minActionableScore {
		return ActivityIdle, npc.CurrentPoi, "none", maxScore
	}

	switch winner {
	case "Hunger":
		activity, location = pickMealForHour(hour), npc.EatingPoi
	case "Fatigue":
		activity, location = ActivitySleeping, npc.HomePoi
	case "Schedule":
		action, _ := npc.ActiveScheduleAt(hour)
		activity, location = action.Activity, action.PoiId
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

func (m *Manager) getNpcsInLocation(blockId string, poi string, npcId string) []*Npc {
	npcs := make([]*Npc, 0)
	for _, npc := range m.npcs {
		if npc.CurrentBlock == blockId && npc.CurrentPoi == poi && npc.Id != npcId {
			npcs = append(npcs, npc)
		}
	}
	return npcs
}
