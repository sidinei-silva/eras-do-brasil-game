package state

import "github.com/sidinei-silva/eras-do-brasil-game/server/internal/domain/world"

type GameState struct {
	// Metadata do tick
	TickCount int64

	// Domínio: World
	GameTime *world.GameTime
	// Blocks    map[string]*Block    // futuro

	// Domínio: NPC
	// NPCs      map[string]*NPC     // futuro

	// Domínio: Combat
	// Combats   map[string]*Combat  // futuro

	// Domínio: Economy
	// Prices    map[string]float64  // futuro
}

func NewGameState() *GameState {
	return &GameState{
		TickCount: 0,
		GameTime:  &world.GameTime{},
		// Blocks:    make(map[string]*Block),
		// NPCs:      make(map[string]*NPC),
		// Combats:   make(map[string]*Combat),
		// Prices:    make(map[string]float64),
	}
}

func (gs *GameState) Snapshot() *GameSnapshot {
	snapshot := &GameSnapshot{
		Tick:     gs.TickCount,
		GameTime: *gs.GameTime,
		// Period:   gs.GameTime.PeriodOfDay,
		// NPCStates: make([]NPCState, 0, len(gs.NPCs)),
		// Combats:   make([]CombatState, 0, len(gs.Combats)),
		// Online:    len(gs.OnlinePlayers),
	}

	// Quando tiveres mapas de NPCs/Players, precisas de iterar e copiar os valores aqui.
	// Exemplo:
	// snap.NPCs = make(map[string]npc.NPC, len(gs.NPCs))
	// for k, v := range gs.NPCs { snap.NPCs[k] = *v } // <- Copia o valor desreferenciado

	return snapshot
}
