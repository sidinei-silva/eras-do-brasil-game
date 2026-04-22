package snapshot

import (
	"github.com/sidinei-silva/eras-do-brasil-game/server/npc"
	"github.com/sidinei-silva/eras-do-brasil-game/server/world"
)

// Snapshot é uma cópia imutável do estado agregado do mundo num tick
// específico. Produzido por Build() e consumido por:
//   - admin hub (serializa e envia ao cliente admin)
//   - persist (futuro, Fase 1+)
//
// Campos devem ser cópias por valor — nunca ponteiros pros structs
// internos dos managers. Isso garante que o leitor não consegue mutar
// estado de jogo por acidente.
type Snapshot struct {
	Tick     int64
	GameTime world.GameTime
	NPCs     []npc.Npc

	// Futuros campos conforme novos managers entrarem:
	// Mobs     map[string]mob.Mob
	// Combats  []combat.Combat
	// Online   int
}

// Build monta uma Snapshot a partir dos managers atuais.
// Chamado pelo GameLoop ao final de cada tick, depois que todos os
// ProcessTick rodaram.
//
// Conforme novos managers forem adicionados, eles entram como parâmetro
// aqui (não em um GameState central).
func Build(tick int64, worldMgr *world.Manager, npcMgr *npc.Manager) *Snapshot {
	snap := &Snapshot{
		Tick:     tick,
		GameTime: worldMgr.GameTime(),
		NPCs:     make([]npc.Npc, 0, len(npcMgr.GetAllNpcs())),
	}

	// Cópia por valor (não ponteiro) para garantir imutabilidade
	// do snapshot após Build retornar.
	for _, v := range npcMgr.GetAllNpcs() {
		snap.NPCs = append(snap.NPCs, *v)
	}

	return snap
}
