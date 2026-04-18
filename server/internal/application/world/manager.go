package world

import (
	"time"

	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/domain/state"
)

type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) ProcessTick(gameState *state.GameState) {
	gameState.GameTime.AdvanceTime(1 * time.Hour)
	// slog.Info("Game time advanced", "tickCount", gameState.TickCount, "currentTime", gameState.GameTime.Time.Format("15:04"), "periodOfDay", gameState.GameTime.PeriodOfDay)
}
