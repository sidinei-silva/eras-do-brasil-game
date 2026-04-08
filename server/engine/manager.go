package engine

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/sidinei-silva/eras-do-brasil-game/server/world"
)

func Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	worldManager := world.NewWorldManager()

	gameLoopReactions := []func(gl *GameLoop){
		func(gl *GameLoop) {
			slog.Info("GameTick", "tickCount", gl.TickCount.Load(), "tickDuration", gl.LastTickDuration)
		},
		func(gl *GameLoop) {
			worldManager.ProcessTick(ctx)
		},
	}

	gameLoop := NewGameLoop(1*time.Second, gameLoopReactions)

	go func() {
		defer wg.Done()
		gameLoop.StartGameLoop(ctx)
	}()
}
