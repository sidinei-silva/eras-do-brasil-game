package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/application/command"
	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/application/world"
	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/domain/state"
	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/infrastructure/engine"
	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/infrastructure/net"
	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/infrastructure/net/api"
	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/infrastructure/net/socket"
)

func main() {

	// Cria um contexto que escuta os sinais SIGINT (Ctrl+C) e SIGTERM
	// 1. Escuta o sinal de parada do Sistema Operacional
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Cria o orquestrador de encerramento
	var wg sync.WaitGroup

	// Cria os objetos
	gameState := state.NewGameState()
	worldManager := world.NewManager()
	cmdManager := command.NewManager()
	cmdQueue := engine.NewCommandQueue()
	router := net.NewCommandRouter(cmdQueue)
	playerHub := socket.NewHub(cmdQueue, router)
	// adminHub := admin.NewHub()

	reactions := func() {
		pendingCommands := cmdQueue.Drain()

		if len(pendingCommands) > 0 {
			cmdManager.ProcessCommands(gameState, pendingCommands)
		}

		worldManager.ProcessTick(gameState)
		// func() {
		// 	snapshot := gameState.Snapshot()
		// 	adminHub.Publish(snapshot)
		// }()
	}

	gameLoop := engine.NewGameLoop(1*time.Second, reactions)

	// Starta goroutines
	mux := http.NewServeMux()
	httpServer := api.NewHTTPServer(mux)
	playerHub.ServeWS(mux, ctx, &wg)
	// adminHub.ServeWS(mux)

	go gameLoop.StartGameLoop(gameState, ctx, &wg)
	go playerHub.Run(ctx, &wg)
	// go adminHub.Run(ctx, &wg)
	go httpServer.StartHTTPServer(ctx, &wg)

	fmt.Println("Sistema rodando. Pressione Ctrl+C para interromper.")

	// A main goroutine bloqueia aqui até que o sistema receba o sinal
	// 2. Trava a execução aqui até o usuário apertar Ctrl+C
	<-ctx.Done()
	fmt.Println("\nSinal de interrupção recebido. Iniciando encerramento gracioso (Graceful Shutdown)...")
	httpServer.StopHTTPServer(ctx)

	// Main bloqueia novamente aqui. Só avança quando todos os `wg.Done()` forem chamados.
	wg.Wait()

	fmt.Println("Todos os processos foram encerrados. Saindo do programa.")

}
