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
	// Load NPCs
	npcs, err := LoadNpcsFromFile()

	if err != nil {
		slog.Error("Erro ao carregar NPCs", "err", err)
		return nil, err
	}

	return &Manager{npcs: npcs}, nil
}

// ProcessTick é o coração do comportamento dos NPCs. É chamado pelo
// game loop a cada tick. Para cada NPC, executa o ciclo:
//  1. Aplica decay contínuo (fome, fadiga, solidão sobem)
//  2. Verifica se atividade discreta atual terminou; se sim, aplica
//     efeito e decide próxima atividade
//  3. Em atividades contínuas, re-avalia decisão
//  4. Garante invariantes (clamp)
func (m *Manager) ProcessTick(gameTime world.GameTime, tickDuration time.Duration) {
	tickHours := tickDuration.Hours()

	for _, npc := range m.npcs {
		// Passo 1: decay sempre roda
		npc.ApplyDecay(tickHours)

		// Passo 2 e 3: lógica de transição de atividade
		if npc.IsDiscreteActivity() {
			if npc.IsActivityComplete(gameTime) {
				npc.ApplyActivityEffects()
				m.decideNextActivity(npc, gameTime)
			}
			// Se não terminou, segue na atividade
		} else {
			// Atividades contínuas podem ser revisitadas a cada tick
			m.decideNextActivity(npc, gameTime)
		}

		// Passo 4: invariantes
		npc.ClampNeeds()
	}
}

// decideNextActivity é o método do Manager que escolhe qual atividade
// o NPC deve fazer agora. Por enquanto usa lógica simples baseada em
// thresholds de need. Será substituído por scoring com pesos na próxima
// rodada (1.3 completo).
func (m *Manager) decideNextActivity(npc *Npc, gameTime world.GameTime) {
	const (
		hungerCritical  = 80.0
		fatigueCritical = 80.0
	)

	hour := gameTime.Time.Hour()
	desiredActivity, desiredLocation := m.computeDesiredActivity(npc, hour)

	if desiredActivity == npc.CurrentActivity && desiredLocation == npc.CurrentZone {
		return
	}

	npc.TransitionTo(desiredActivity, desiredLocation, gameTime)

}

func (m *Manager) computeDesiredActivity(npc *Npc, hour int) (Activity, string) {
	const (
		hungerCritical  = 80.0
		fatigueCritical = 80.0
	)

	if npc.Needs.Hunger >= hungerCritical {
		return pickMealForHour(hour), npc.EatingLocation
	}
	if npc.Needs.Fatigue >= fatigueCritical {
		return ActivitySleeping, npc.HomeLocation
	}
	if action, found := npc.ActiveScheduleAt(hour); found {
		return action.Activity, action.Location
	}
	return ActivityIdle, npc.CurrentZone
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
