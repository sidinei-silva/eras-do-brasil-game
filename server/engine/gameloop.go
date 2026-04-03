package engine

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sidinei-silva/eras-do-brasil-game/server/world"
)

type Status struct {
	Running          bool          `json:"running"`
	Interval         time.Duration `json:"interval"`
	TickCount        int64         `json:"tick_count"`
	StartedAt        time.Time     `json:"started_at"`
	UpTime           time.Duration `json:"up_time"`
	LastTickDuration time.Duration `json:"last_tick_duration"`
}

type SnapshotPublisher interface {
	Publish(snap world.Snapshot)
}

type GameLoop struct {
	interval         time.Duration
	world            *world.World
	publisher        SnapshotPublisher
	running          atomic.Bool
	tickCount        atomic.Uint64
	startedAt        time.Time
	cancel           context.CancelFunc
	mu               sync.RWMutex
	lastTickDuration time.Duration
}

func NewGameLoop(interval time.Duration, world *world.World, publisher SnapshotPublisher) *GameLoop {
	return &GameLoop{
		interval:  interval,
		world:     world,
		publisher: publisher,
	}
}

func (g *GameLoop) Start(ctx context.Context) {
	if !g.running.CompareAndSwap(false, true) {
		slog.Warn("Game loop is already running")
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	g.cancel = cancel
	g.startedAt = time.Now()
	slog.Info("Game loop started")

	ticker := time.NewTicker(g.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			g.running.Store(false)
			slog.Info("Game loop cancelled")
			return

		case <-ticker.C:
			start := time.Now()
			snap := g.world.ProcessTick()

			if g.publisher != nil {
				g.publisher.Publish(snap)
			}

			g.tickCount.Add(1)
			tickDuration := time.Since(start)

			g.mu.Lock()
			g.lastTickDuration = tickDuration
			g.mu.Unlock()
		}
	}
}

func (g *GameLoop) Stop() {
	if !g.running.CompareAndSwap(true, false) {
		slog.Warn("Game loop is not running")
		return
	}

	if g.cancel != nil {
		g.cancel()
	}

	slog.Info("Game loop stopped")
}

func (g *GameLoop) SetPublisher(p SnapshotPublisher) {
	g.publisher = p
}

func (g *GameLoop) Status() Status {
	g.mu.RLock()
	last := g.lastTickDuration
	g.mu.RUnlock()

	started := g.startedAt
	uptime := time.Duration(0)
	if !started.IsZero() {
		uptime = time.Since(started)
	}

	return Status{
		Running:          g.running.Load(),
		Interval:         g.interval,
		TickCount:        int64(g.tickCount.Load()),
		StartedAt:        started,
		UpTime:           uptime,
		LastTickDuration: last,
	}
}
