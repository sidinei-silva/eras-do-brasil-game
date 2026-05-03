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

	charmlog "github.com/charmbracelet/log"
	"github.com/sidinei-silva/eras-do-brasil-game/server/command"
	"github.com/sidinei-silva/eras-do-brasil-game/server/engine"
	"github.com/sidinei-silva/eras-do-brasil-game/server/net/api"
	adminSocket "github.com/sidinei-silva/eras-do-brasil-game/server/net/socket/admin"
	playerSocket "github.com/sidinei-silva/eras-do-brasil-game/server/net/socket/player"
	"github.com/sidinei-silva/eras-do-brasil-game/server/npc"
	"github.com/sidinei-silva/eras-do-brasil-game/server/snapshot"
	"github.com/sidinei-silva/eras-do-brasil-game/server/world"
)

func main() {

	charmLogger := charmlog.NewWithOptions(os.Stderr, charmlog.Options{
		Level:           charmlog.DebugLevel,
		TimeFormat:      time.TimeOnly,
		ReportCaller:    false,
		ReportTimestamp: true,
	})
	slog.SetDefault(slog.New(charmLogger))

	// Cria um contexto que escuta os sinais SIGINT (Ctrl+C) e SIGTERM
	// 1. Escuta o sinal de parada do Sistema Operacional
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Cria o orquestrador de encerramento
	var wg sync.WaitGroup

	// Cria os objetos
	// 1. Instanciação (Injeção de Dependência)
	cmdQueue := command.NewCommandQueue()

	// Hubs
	playerHub := playerSocket.NewHub(cmdQueue)
	adminHub := adminSocket.NewHub()

	snapManager := snapshot.NewManager()

	// Roteadores
	playerRouter := command.NewPlayerRouter(cmdQueue)
	adminRouter := command.NewAdminRouter(cmdQueue, adminHub, snapManager)

	//Injeta o Router no Hub
	adminHub.SetRouter(adminRouter)
	playerHub.SetRouter(playerRouter)

	// tickDuration: quanto tempo de jogo avança por tick (tickInterval = 1s).
	// Produção: 12 * time.Second → 1h de jogo = 5 min reais, 1 dia de jogo = 2h reais (GDD §8.2 e §8.9)
	// Atual: 1h para observar o comportamento dos NPCs com o passar do tempo de forma acelerada.
	tickDuration := 1 * time.Hour
	worldManager, worldErr := world.NewManager(tickDuration)
	if worldErr != nil {
		slog.Error("Erro ao criar World Manager:", "worldErr", worldErr)
		os.Exit(1)
	}

	slog.Info("Game Time inicializado", "gameTime", worldManager.GameTime())
	blocks := worldManager.GetAllBlocks()
	slog.Info("Blocos carregados no World Manager", "count", len(blocks))

	npcManager, npcErr := npc.NewManager()
	if npcErr != nil {
		slog.Error("Erro ao criar NPC Manager:", "npcErr", npcErr)
		os.Exit(1)
	}

	npcs := npcManager.GetAllNpcs()
	slog.Info("NPCs carregados no NPC Manager", "count", len(npcs))

	gameLoop := engine.NewGameLoop()

	// 2. Definição das Reações do Tick
	gameLoop.SetReactionsForTick(func() {
		// Passo 1: Processa comandos (Player + Admin)
		pendingCommands := cmdQueue.Drain()

		if len(pendingCommands) > 0 {
			command.ProcessPlayerCommands(gameLoop.TickCount(), pendingCommands)
		}
		// Passo 2: Evolui o mundo
		worldManager.ProcessTick()

		// Passo 3: Atualiza os NPCs
		npcManager.ProcessTick(worldManager.GameTime(), worldManager.TickDuration())

		// Passo 4: Modo Deus - Tira a foto e manda para o Admin Hub
		snapshot := snapManager.BuildSnapshot(gameLoop.TickCount(), worldManager, npcManager)
		adminHub.Publish(snapshot)
	})

	// Starta goroutines
	// 3. Inicialização das Redes e Processos
	mux := http.NewServeMux()
	httpServer := api.NewHTTPServer(mux)

	playerHub.ServeWS(mux, ctx, &wg)
	adminHub.ServeWS(mux, ctx, &wg)

	// Roda tudo em paralelo
	go gameLoop.StartGameLoop(ctx, &wg)
	go playerHub.Run(ctx, &wg)
	go adminHub.Run(ctx, &wg)
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
