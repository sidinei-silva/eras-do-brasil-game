package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sidinei-silva/eras-do-brasil-game/server/engine"
	handlers "github.com/sidinei-silva/eras-do-brasil-game/server/http"
	"github.com/sidinei-silva/eras-do-brasil-game/server/world"
)

func startGame(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	gameTime := world.NewGameTime()

	gameLoopReactions := []func(gl *engine.GameLoop){
		func(gl *engine.GameLoop) {
			slog.Info("GameTick", "tickCount", gl.TickCount.Load(), "tickDuration", gl.LastTickDuration)
		},

		// Advance game time every tick
		func(gl *engine.GameLoop) {
			gameTime.AdvanceTime(1 * time.Minute)
			slog.Info("Game time advanced", "tickCount", gl.TickCount.Load(), "currentTime", gameTime.Time.Format("15:04"), "periodOfDay", gameTime.PeriodOfDay)
		},
	}

	gameLoop := engine.NewGameLoop(1*time.Second, gameLoopReactions)

	go func() {
		defer wg.Done()
		gameLoop.StartGameLoop(ctx)
	}()

}

func startServer(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	mux := http.NewServeMux()
	handlers.Handlers(mux)

	server := &http.Server{
		Addr:        ":8080",
		ReadTimeout: 5 * time.Second,
		Handler:     mux,
	}

	// 1. Goroutine para rodar o servidor (escutar requisições)
	go func() {
		slog.Info("server iniciado", "addr", server.Addr)
		slog.Info("admin dashboard disponível em http://localhost:8080/admin/")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("erro no servidor", "err", err)
		}
	}()

	// 2. Goroutine para escutar o cancelamento e fazer o Shutdown
	go func() {
		defer wg.Done() // Avisa a main que o desligamento do servidor terminou

		// Trava AQUI até o Ctrl+C acontecer
		<-ctx.Done()

		slog.Info("Iniciando shutdown do servidor...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("erro durante o shutdown ou timeout atingido", "err", err)
		} else {
			slog.Info("Servidor encerrado com sucesso.")
		}
	}()

}

func main() {
	// Cria um contexto que escuta os sinais SIGINT (Ctrl+C) e SIGTERM
	// 1. Escuta o sinal de parada do Sistema Operacional
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Cria o orquestrador de encerramento
	var wg sync.WaitGroup

	startGame(ctx, &wg)
	startServer(ctx, &wg)

	fmt.Println("Sistema rodando. Pressione Ctrl+C para interromper.")

	// A main goroutine bloqueia aqui até que o sistema receba o sinal
	// 2. Trava a execução aqui até o usuário apertar Ctrl+C
	<-ctx.Done()
	fmt.Println("\nSinal de interrupção recebido. Iniciando encerramento gracioso (Graceful Shutdown)...")

	// Main bloqueia novamente aqui. Só avança quando todos os `wg.Done()` forem chamados.
	wg.Wait()
	
	fmt.Println("Todos os processos foram encerrados. Saindo do programa.")

}
