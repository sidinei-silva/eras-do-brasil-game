package world

import "time"

type Manager struct {
	gameTime *GameTime
	// futuramente: blocks, climate
}

func NewManager() *Manager {
	return &Manager{gameTime: &GameTime{}}
}

func (m *Manager) ProcessTick() {
	m.gameTime.AdvanceTime(1 * time.Hour)
}

func (m *Manager) GameTime() GameTime {
	return *m.gameTime // retorna cópia, leitor não pode mutar
}
