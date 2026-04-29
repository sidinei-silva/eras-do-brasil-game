package engine

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

// tickInterval é o tempo real entre cada tick do loop principal.
// Mantido em 1s para garantir responsividade do jogador (comandos processados em até 1s).
// A velocidade do tempo de jogo é controlada pelo tickDuration no WorldManager, não aqui.
const tickInterval = 1 * time.Second

type GameLoop struct {
	// tickCount armazena a quantidade total de ticks já executados no loop.
	tickCount atomic.Int64
	// running indica se o loop principal está em execução no momento.
	running atomic.Bool
	// cancel guarda a função usada para cancelar o contexto interno do loop.
	cancel context.CancelFunc
	// LastTickDuration registra quanto tempo o último tick levou para ser processado.
	LastTickDuration time.Duration
	// reactionsForTick contém a função chamada a cada tick para processar as reações do jogo.
	reactionsForTick func()
}

func NewGameLoop() *GameLoop {
	return &GameLoop{
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

	ticker := time.NewTicker(tickInterval)
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
