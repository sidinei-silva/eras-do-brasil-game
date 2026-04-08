package engine

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"
)

type GameLoop struct {
	interval         time.Duration
	running          atomic.Bool
	TickCount        atomic.Int64
	cancel           context.CancelFunc
	LastTickDuration time.Duration
	reactionsForTick []func(gl *GameLoop)
}

func NewGameLoop(interval time.Duration, reactions []func(gl *GameLoop)) *GameLoop {
	return &GameLoop{
		interval:         interval,
		reactionsForTick: reactions,
	}
}

func (gl *GameLoop) StartGameLoop(ctx context.Context) {
	if !gl.running.CompareAndSwap(false, true) {
		slog.Warn("Game loop is already running.")
		return
	}

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
			gl.TickCount.Add(1)
			for _, reaction := range gl.reactionsForTick {
				reaction(gl)
			}
			gl.LastTickDuration = time.Since(start)
		}
	}
}

func (gl *GameLoop) StopGameLoop() {
	if !gl.running.CompareAndSwap(true, false) {
		slog.Warn("Game Loop is not running")
		return
	}

	slog.Info("Game loop stopped")
}
