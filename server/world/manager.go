package world

import (
	"log/slog"
	"time"
)

type Manager struct {
	gameTime     *GameTime     // Gerencia o tempo do jogo
	blocks       []*Block      // Blocos do mundo
	tickDuration time.Duration // Quanto tempo o jogo avança a cada tick
}

func NewManager(tickDuration time.Duration) (*Manager, error) {

	// Carrega blocos
	blocks, err := LoadBlocksFromFile()
	gameTime := NewGameTime()

	if err != nil {
		slog.Error("Erro ao carregar blocos", "err", err)
		return nil, err
	}

	return &Manager{gameTime: &gameTime, blocks: blocks, tickDuration: tickDuration}, nil
}

func (m *Manager) ProcessTick() {
	m.gameTime.AdvanceTime(m.tickDuration)
}

func (m *Manager) GameTime() GameTime {
	return *m.gameTime // retorna cópia, leitor não pode mutar
}

func (m *Manager) GetBlockById(id string) (*Block, bool) {
	for _, block := range m.blocks {
		if block.Id == id {
			return block, true
		}
	}
	return nil, false
}

func (m *Manager) GetAllBlocks() []*Block {
	blocks := make([]*Block, 0, len(m.blocks))
	for _, block := range m.blocks {
		blocks = append(blocks, block)
	}
	return blocks
}

func (m *Manager) TickDuration() time.Duration {
	return m.tickDuration
}
