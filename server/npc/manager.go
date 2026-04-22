package npc

import (
	"log/slog"

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

func (m *Manager) ProcessTick(gameTime world.GameTime) {
	// Aqui é onde a lógica de comportamento dos NPCs seria processada a cada tick.
	// Por exemplo, você poderia iterar sobre os NPCs e atualizar suas necessidades, mudar suas atividades com base no horário do jogo, etc.
	for _, npc := range m.npcs {
		npc.UpdateNeeds(gameTime)
		npc.CurrentActivityAndLocation(gameTime)
	}
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
