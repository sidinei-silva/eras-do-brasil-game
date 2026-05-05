package snapshot

import (
	"log/slog"

	"github.com/sidinei-silva/eras-do-brasil-game/server/npc"
	"github.com/sidinei-silva/eras-do-brasil-game/server/world"
)

// Snapshot é uma cópia imutável do estado agregado do mundo num tick
// específico. Produzido por Build() e consumido por:
//   - hub admin (serializa e envia ao cliente admin)
//   - persistência (futuro, Fase 1+)
//
// Campos devem ser cópias por valor — nunca ponteiros pros structs
// internos dos managers. Isso garante que o leitor não consegue mutar
// estado de jogo por acidente.
type Snapshot struct {
	Tick     int64
	GameTime world.GameTime
	NPCs     []npc.Npc
	Blocks   []world.Block

	// Futuros campos conforme novos managers entrarem:
	// Mobs     map[string]mob.Mob
	// Combats  []combat.Combat
	// Online   int
}

func build(tick int64, worldMgr *world.Manager, npcMgr *npc.Manager) *Snapshot {
	snap := &Snapshot{
		Tick:     tick,
		GameTime: worldMgr.GameTime(),
		NPCs:     make([]npc.Npc, 0, len(npcMgr.GetAllNpcs())),
		Blocks:   make([]world.Block, 0, len(worldMgr.GetAllBlocks())),
	}

	// Cópia por valor (não ponteiro) para garantir imutabilidade
	// do Snapshot após Build retornar.
	for _, v := range npcMgr.GetAllNpcs() {
		snap.NPCs = append(snap.NPCs, *v)
	}

	for _, v := range worldMgr.GetAllBlocks() {
		snap.Blocks = append(snap.Blocks, *v)
	}

	return snap
}

func (s *Snapshot) GetNPCById(id string) (*npc.Npc, bool) {
	for _, npc := range s.NPCs {
		if npc.Id == id {
			return &npc, true
		}
	}
	return nil, false
}

func (s *Snapshot) GetNpcScores(npcId string) map[string]float64 {
	npc, found := s.GetNPCById(npcId)

	if !found {
		slog.Warn("NPC não encontrado para calcular scores", "npc_id", npcId)
		return nil
	}
	gameTime := s.GetGameTime()
	scores := npc.CalculateScores(gameTime.Time.Hour())
	return scores
}

func (s *Snapshot) GetBlockById(id string) (*world.Block, bool) {
	for _, block := range s.Blocks {
		if block.Id == id {
			return &block, true
		}
	}
	return nil, false
}

func (s *Snapshot) GetNPCsInZone(zoneId string) []npc.Npc {
	npcsInZone := make([]npc.Npc, 0)
	for _, npc := range s.NPCs {
		if npc.CurrentBlock == zoneId {
			npcsInZone = append(npcsInZone, npc)
		}
	}
	return npcsInZone
}

func (s *Snapshot) GetAllNPCs() []npc.Npc {
	return s.NPCs
}

func (s *Snapshot) GetAllBlocks() []world.Block {
	return s.Blocks
}

func (s *Snapshot) GetGameTime() world.GameTime {
	return s.GameTime
}
