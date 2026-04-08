package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sidinei-silva/eras-do-brasil-game/server/engine"
	"github.com/sidinei-silva/eras-do-brasil-game/server/net"
)

func main() {

	// Cria um contexto que escuta os sinais SIGINT (Ctrl+C) e SIGTERM
	// 1. Escuta o sinal de parada do Sistema Operacional
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Cria o orquestrador de encerramento
	var wg sync.WaitGroup

	engine.Start(ctx, &wg)
	net.Start(ctx, &wg)

	fmt.Println("Sistema rodando. Pressione Ctrl+C para interromper.")

	// A main goroutine bloqueia aqui até que o sistema receba o sinal
	// 2. Trava a execução aqui até o usuário apertar Ctrl+C
	<-ctx.Done()
	fmt.Println("\nSinal de interrupção recebido. Iniciando encerramento gracioso (Graceful Shutdown)...")

	// Main bloqueia novamente aqui. Só avança quando todos os `wg.Done()` forem chamados.
	wg.Wait()

	fmt.Println("Todos os processos foram encerrados. Saindo do programa.")

}
