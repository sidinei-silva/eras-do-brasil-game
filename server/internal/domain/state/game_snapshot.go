package state

import "github.com/sidinei-silva/eras-do-brasil-game/server/internal/domain/world"

type GameSnapshot struct {
	Tick     int64
	GameTime world.GameTime // "15:04" formatado
	// Period   string         // "Manhã", "Tarde", etc.
	// NPCStates []NPCState  // futuro: cópia flat, não ponteiros
	// Combats   []CombatState
	// Online    int
}
