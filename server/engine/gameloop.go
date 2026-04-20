package engine

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

type GameLoop struct {
	tickCount        atomic.Int64
	interval         time.Duration
	running          atomic.Bool
	cancel           context.CancelFunc
	LastTickDuration time.Duration
	reactionsForTick func()
}

func NewGameLoop(interval time.Duration) *GameLoop {
	return &GameLoop{
		interval:         interval,
		reactionsForTick: func() {},
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

func (gl *GameLoop) TickCount() int64 {
	return gl.tickCount.Load()
}

func (gl *GameLoop) StartGameLoop(ctx context.Context, wg *sync.WaitGroup) {
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
			gl.tickCount.Add(1)
			gl.reactionsForTick()
			gl.LastTickDuration = time.Since(start)
		}
	}
}

func (gl *GameLoop) SetReactionsForTick(reactions func()) {
	gl.reactionsForTick = reactions
}
