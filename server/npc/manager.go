package npc

import (
	"log/slog"

	"github.com/sidinei-silva/eras-do-brasil-game/server/world"
)

type Manager struct {
	npcs map[string]*Npc
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

func (m *Manager) All() map[string]*Npc { return m.npcs }
func (m *Manager) Get(id string) (*Npc, bool) {
	n, ok := m.npcs[id]
	return n, ok
}
