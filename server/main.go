package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sidinei-silva/eras-do-brasil-game/server/engine"
	"github.com/sidinei-silva/eras-do-brasil-game/server/socket"
	"github.com/sidinei-silva/eras-do-brasil-game/server/world"
)

// ServerError é um tipo de erro customizado (não obrigatório, mas bom para debug).
// Implementar Error() torna ServerError um "error" válido em Go.
type ServerError struct {
	Message string
	Err     error
}

// Error() faz ServerError implementar a interface "error".
// Interface = contrato: qualquer tipo que tiver Error() é considerado erro.
func (se *ServerError) Error() string {
	return se.Message
}

// Unwrap() permite "desembrulhar" erros (Go 1.13+), útil para debugging.
func (se *ServerError) Unwrap() error {
	return se.Err
}

// main() é a função de entrada do programa.
// Quando você executa `go run main.go`, Go procura por uma função main() em package main.
func main() {
	// signal.NotifyContext(parent_ctx, sinais...) retorna:
	// 1. Um novo context que é cancelado quando um dos sinais é recebido
	// 2. Uma função stop() que pode cancelar o context manualmente
	//
	// syscall.SIGINT = Ctrl+C
	// syscall.SIGTERM = sinal do SO para "termine o processo"
	//
	// Portanto: context.Background() é "contexto raiz", que se cancela se o SO mandar parar.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// defer stop() = garante que a função stop é chamada ao final,
	// mesmo que haja erro ou panic. Não é obrigatório aqui, mas é bom prática.
	defer stop()

	// Cria uma nova instância do mundo do jogo.
	// NewWorld() retorna um PONTEIRO (*World), então mundo é um endereço de memória.
	// Todas as mudanças ao mundo afetam o original, não uma cópia.
	mundo := world.NewWorld()

	// Cria o game loop com intervalo de 500 milissegundos entre cada tick.
	// time.Millisecond = constante que Go oferece (1 ms = 0,001 seg).
	// Então 500 * time.Millisecond = 500 ms = 0,5 segundos entre ticks.
	loop := engine.NewGameLoop(500*time.Millisecond, mundo)

	// go loop.Start(ctx) = cria uma GOROUTINE (thread leve) que roda loop.Start().
	// "go" = execute isto em paralelo, não bloqueie aqui.
	// Se não tivesse "go", main() esperaria que loop.Start() terminasse (nunca, é infinito).
	// Com "go", main() continua e pode fazer outras coisas.
	//
	// ctx é o contexto que se cancela se Ctrl+C ou SIGTERM.
	// Quando ctx é cancelado, loop.Start() vê <-ctx.Done() e sai.
	go loop.Start(ctx)

	hub := socket.NewHub()
	go hub.Run(ctx)

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				hub.Broadcast(socket.OutboundEvent{
					Type: "world_snapshot",
					Data: mundo.Snapshot(),
				})
			}
		}
	}()

	// Define um endpoint HTTP que retorna status em JSON.
	// http.HandleFunc(caminho, função_handler) = "quando alguém acessar caminho, chame a função".
	http.HandleFunc("/admin/status", func(w http.ResponseWriter, r *http.Request) {
		// w = "resposta" que vai voltar para o cliente
		// r = "requisição" que veio do cliente

		// w.Header().Set(...) = adiciona um header HTTP na resposta.
		// "Content-Type: application/json" diz ao navegador "isto é JSON, não HTML".
		w.Header().Set("Content-Type", "application/json")

		// json.NewEncoder(w).Encode(...) = converte uma struct Go em JSON e envia via w.
		// map[string]any = um "mapa" onde chave é string, valor é "any" (qualquer tipo).
		// Isso permite misturar diferentes tipos na resposta JSON.
		_ = json.NewEncoder(w).Encode(map[string]any{
			// "game_loop": status atual do loop (rodando?, quantos ticks?, uptime?)
			"game_loop": loop.Status(),
			// "world": fotografia do estado do mundo (hora do jogo, período do dia)
			// IMPORTANTE: Snapshot() SÓ LÊ, não modifica. ProcessTick() MODIFICA!
			"world": mundo.Snapshot(),
			"lobby": map[string]any{
				"online": hub.OnlineCount(),
			},
		})
		// O "_" descarta o erro se houver (não é bom prática, mas por enquanto é OK para debug).
	})

	http.HandleFunc("/ws", socket.WsHandler(hub))

	// http.Server = configuração do servidor HTTP.
	// &http.Server{} = cria uma nova instância (& = endereço).
	server := &http.Server{
		Addr:         ":8080",         // porta 8080 (http://localhost:8080)
		ReadTimeout:  5 * time.Second, // máximo 5 segundos para cliente enviar requisição completa
		WriteTimeout: 5 * time.Second, // máximo 5 segundos para servidor enviar resposta completa
	}

	// go func() { ... }() = cria uma goroutine com uma função anônima (sem nome).
	// "Anônima" = função definida inline, só usada aqui.
	// Este padrão é comum para executar algo em paralelo sem quebrar o fluxo.
	go func() {
		// slog.Info(...) = log no console com nível INFO (informação útil).
		slog.Info("server iniciado", "addr", server.Addr)

		// server.ListenAndServe() = bloqueia aqui, aguardando requisições HTTP.
		// Se houver erro e não for "servidor fechou", avisa.
		// err != http.ErrServerClosed = ignora o erro esperado quando Stop() é chamado.
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("erro no servidor", "err", err)
			// stop() = cancela o context, que avisa o loop para parar.
			stop()
		}
	}()

	// <-ctx.Done() = bloqueia aqui até ctx ser cancelado.
	// ctx é cancelado quando:
	// 1. Ctrl+C for pressionado (SIGINT)
	// 2. O SO enviar SIGTERM
	// Então este programa "dorme" aqui enquanto está rodando.
	<-ctx.Done()

	// Código abaixo só executa após ctx.Done() ser ativado (Ctrl+C ou SIGTERM).

	// Para o game loop graciosamente.
	// loop.Stop() chama cancel(), que ativa <-ctx.Done() dentro de loop.Start().
	loop.Stop()

	// context.WithTimeout(...) = cria um novo context que se cancela após 5 segundos.
	// Usado para garantir que o Shutdown não espera para sempre.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel() = cancela esse context extra quando main sair (limpeza).
	defer cancel()

	// server.Shutdown(shutdownCtx) = desliga o servidor HTTP graciosamente.
	// Termina conexões em andamento e para de aceitar novas requisições.
	// Se demorar mais que 5 segundos, shutdownCtx se cancela e força saída.
	if err := server.Shutdown(shutdownCtx); err != nil {
		// Se erro, avisa no console.
		slog.Error("erro no shutdown", "err", err)
		// os.Exit(1) = fecha o programa com código de erro 1.
		// Código 0 = sucesso, qualquer outro = erro (para ferramentas externas saberem).
		os.Exit(1)
	}

	// Se chegou aqui, desligamento foi bem-sucedido.
	slog.Info("programa finalizado")
}
