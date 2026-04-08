package world

import (
	"context"
	"log/slog"
	"time"
)

type WorldManager struct {
	gameTime GameTime
}

func NewWorldManager() *WorldManager {
	gameTime := NewGameTime()

	return &WorldManager{
		gameTime: gameTime,
	}
}

func (wm *WorldManager) ProcessTick(ctx context.Context) {
	// Achar uma forma de ter um GameState compartilhado para todos os managers, e pegar o tickCount de lá
	tickCount := "unknown"

	wm.gameTime.AdvanceTime(1 * time.Minute)
	slog.Info("Game time advanced", "tickCount", tickCount, "currentTime", wm.gameTime.Time.Format("15:04"), "periodOfDay", wm.gameTime.PeriodOfDay)
}
