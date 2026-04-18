package npc

import (
	"fmt"

	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/domain/state"
	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/infrastructure/data"
)

type Manager struct{}

func NewManager(gameState *state.GameState) (*Manager, error) {
	// Load NPCs
	npcs, err := data.LoadNpcsFromFile()

	if err != nil {
		fmt.Println("Erro ao carregar NPCs:", err)
		return nil, err
	}

	gameState.NPCs = npcs

	return &Manager{}, nil
}

func (m *Manager) ProcessTick(gameState *state.GameState) {
	// Aqui é onde a lógica de comportamento dos NPCs seria processada a cada tick.
	// Por exemplo, você poderia iterar sobre os NPCs e atualizar suas necessidades, mudar suas atividades com base no horário do jogo, etc.
	gameTime := gameState.GameTime
	for _, npc := range gameState.NPCs {
		npc.UpdateNeeds(*gameTime)
		npc.CurrentActivityAndLocation(*gameTime)
	}
}
