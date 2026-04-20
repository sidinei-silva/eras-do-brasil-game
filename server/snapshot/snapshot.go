package snapshot

import (
	"github.com/sidinei-silva/eras-do-brasil-game/server/npc"
	"github.com/sidinei-silva/eras-do-brasil-game/server/world"
)

type Snapshot struct {
	Tick     int64
	GameTime world.GameTime
	NPCs     map[string]npc.Npc
	// Period   string         // "Manhã", "Tarde", etc.
	// NPCStates []NPCState  // futuro: cópia flat, não ponteiros
	// Combats   []CombatState
	// Online    int
}

func Build(tick int64, worldMgr *world.Manager, npcMgr *npc.Manager) *Snapshot {
	snap := &Snapshot{
		Tick:     tick,
		GameTime: worldMgr.GameTime(),
		NPCs:     make(map[string]npc.Npc, len(npcMgr.All())),
		// Period:   gs.GameTime.PeriodOfDay,
		// NPCStates: make([]NPCState, 0, len(gs.NPCs)),
		// Combats:   make([]CombatState, 0, len(gs.Combats)),
		// Online:    len(gs.OnlinePlayers),
	}

	// Quando tiveres mapas de NPCs/Players, precisas de iterar e copiar os valores aqui.
	// Exemplo:
	// snap.NPCs = make(map[string]npc.NPC, len(gs.NPCs))
	// for k, v := range gs.NPCs { snap.NPCs[k] = *v } // <- Copia o valor desreferenciado

	for k, v := range npcMgr.All() {
		snap.NPCs[k] = *v
	}
	return snap
}
