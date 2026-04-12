package engine

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/domain/state"
)

type GameLoop struct {
	interval         time.Duration
	running          atomic.Bool
	cancel           context.CancelFunc
	LastTickDuration time.Duration
	reactionsForTick func()
}

func NewGameLoop(interval time.Duration, reactions func()) *GameLoop {
	return &GameLoop{
		interval:         interval,
		reactionsForTick: reactions,
	}
}

func (gl *GameLoop) StopGameLoop() {
	if gl.cancel != nil {
		gl.cancel()
	}
	gl.running.Store(false)
}

func (gl *GameLoop) IsRunning() bool {
	return gl.running.Load()
}

func (gl *GameLoop) StartGameLoop(gameState *state.GameState, ctx context.Context, wg *sync.WaitGroup) {
	if !gl.running.CompareAndSwap(false, true) {
		slog.Warn("Game loop is already running.")
		return
	}

	wg.Add(1)
	defer wg.Done()

	ctx, cancel := context.WithCancel(ctx)
	gl.cancel = cancel

	slog.Info("Starting game loop")

	ticker := time.NewTicker(gl.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Context canceled, stopping game loop")
			gl.StopGameLoop()
			return
		case <-ticker.C:
			start := time.Now()
			gameState.TickCount++
			gl.reactionsForTick()
			gl.LastTickDuration = time.Since(start)
		}
	}
}
