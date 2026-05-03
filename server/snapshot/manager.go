package snapshot

import (
	"github.com/sidinei-silva/eras-do-brasil-game/server/npc"
	"github.com/sidinei-silva/eras-do-brasil-game/server/world"
)

type Manager struct {
	snapshot *Snapshot
	// Futuramente, podemos adicionar um canal para notificar o Hub Admin de novas snapshots, ou o Hub pode simplesmente puxar a última snapshot a cada tick.
}

func NewManager() *Manager {

	return &Manager{
		snapshot: nil, // Começa sem snapshot; a primeira será construída no primeiro tick do jogo.
	}
}

// Build monta uma Snapshot a partir dos managers atuais.
// Chamado pelo GameLoop ao final de cada tick, depois que todos os
// ProcessTick rodaram.
//
// Conforme novos managers forem adicionados, eles entram como parâmetro
// aqui (não em um GameState central).
func (m *Manager) BuildSnapshot(tickCount int64, worldManager *world.Manager, npcManager *npc.Manager) *Snapshot {
	snap := build(tickCount, worldManager, npcManager)
	m.snapshot = snap
	return snap
}

func (m *Manager) GetSnapshot() *Snapshot {
	return m.snapshot
}
