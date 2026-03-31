package engine

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sidinei-silva/eras-do-brasil-game/server/world"
)

// Status representa a situação atual do game loop.
// json: tags definem nomes para JSON quando isso for serializado (convertido em JSON para HTTP).
type Status struct {
	Running          bool          `json:"running"`            // verdadeiro se loop está rodando
	Interval         time.Duration `json:"interval"`           // intervalo entre ticks
	TickCount        int64         `json:"tick_count"`         // número total de ticks processados
	StartedAt        time.Time     `json:"started_at"`         // quando o loop começou
	UpTime           time.Duration `json:"up_time"`            // quanto tempo o loop está rodando
	LastTickDuration time.Duration `json:"last_tick_duration"` // quanto tempo o último tick levou
}

// SnapshotPublisher é a interface que o GameLoop usa para publicar o estado pós-tick.
// Qualquer tipo que implemente Publish pode receber snapshots — sem dependência de socket.
type SnapshotPublisher interface {
	Publish(snap world.Snapshot)
}

// GameLoop é o orquestrador central do jogo.
// Ele processa ticks em intervalo fixo e coordena o processamento do mundo.
type GameLoop struct {
	// interval: tempo entre cada tick do jogo (ex: 500ms).
	interval time.Duration

	// world: PONTEIRO (*) para o mundo.
	// * = é um endereço de memória, não uma cópia. Se alterar world, altera o original.
	// Se fosse só `world world.World`, seria uma cópia, e mudanças não afetariam o original.
	world *world.World

	// publisher: destino dos snapshots pós-tick (ex: hub WebSocket).
	// nil = sem publicação (útil em testes).
	publisher SnapshotPublisher

	// running: ATOMIC BOOL = booleano "thread-safe".
	// "thread-safe" = múltiplas goroutines podem ler/escrever sem corromper dados.
	// atomic = operações indivisíveis, garantidas pelo hardware/SO.
	// Por que atomic e não um mutex? Porque é mais rápido para uma única variável.
	running atomic.Bool

	// tickCount: ATOMIC UINT64 = contador "thread-safe" de ticks.
	// uint64 = número inteiro sem sinal, até 18 quintilhões.
	// Usamos atomic porque múltiplas goroutines podem ler esse número.
	tickCount atomic.Uint64

	// startedAt: hora em que o loop foi iniciado.
	// time.Time = tipo que Go oferece para datas/horas.
	startedAt time.Time

	// cancel: FUNÇÃO que cancela o contexto do loop.
	// context.CancelFunc = função que, quando chamada, sinaliza ao loop para parar.
	// É usado para "desligar" o loop graciosamente.
	cancel context.CancelFunc

	// mu: mutex para proteger lastTickDuration contra acesso simultâneo.
	// RWMutex = permite múltiplos leitores de uma vez, MAS só um escritor.
	// Lock/Unlock = escrita exclusiva (ninguém mais acessa).
	// RLock/RUnlock = leitura compartilhada (outros podem ler enquanto isso).
	mu sync.RWMutex

	// lastTickDuration: tempo que o último tick levou para processar.
	// Protegido por mu (o mutex acima) porque pode ser alterado e lido simultaneamente.
	lastTickDuration time.Duration
}

// NewGameLoop cria uma nova instância do game loop.
// Retorna um PONTEIRO (*GameLoop) porque o loop é "mortal" — vive enquanto o programa roda.
// Se retornasse cópia, mudanças ao loop não afetariam outras goroutines.
// publisher recebe o snapshot após cada tick — pode ser nil.
func NewGameLoop(interval time.Duration, world *world.World, publisher SnapshotPublisher) *GameLoop {
	return &GameLoop{
		interval:  interval,
		world:     world,
		publisher: publisher,
	}
}

// Start inicia o game loop em uma goroutine separada.
// ctx é um CONTEXT = mecanismo para sinalizar "cancele o que está fazendo".
// Um context pode vir de um sinal (Ctrl+C), timeout, ou ser cancelado manualmente.
func (g *GameLoop) Start(ctx context.Context) {
	// CompareAndSwap(esperado, novo) = operação ATÔMICA que:
	// 1. verifica se running é FALSE (false = 0)
	// 2. se true, muda para TRUE (true = 1) e retorna true
	// 3. se false não era, não muda nada e retorna false
	// Isso evita que Start seja chamado duas vezes simultaneamente.
	if !g.running.CompareAndSwap(false, true) {
		// Se já estava rodando, avisa e sai.
		slog.Warn("Game loop is already running")
		return
	}

	// context.WithCancel(ctx) = pega um context existente e cria uma VERSÃO CANCELLÁVEL.
	// Retorna 2 coisas:
	// 1. Novo context (que herda cancelamentos do anterior)
	// 2. Uma função cancel() que, quando chamada, cancela rapidinho
	ctx, cancel := context.WithCancel(ctx)

	// Guarda a função cancel para usar depois em Stop().
	g.cancel = cancel

	// Guarda a hora que o loop começou (para calcular uptime depois).
	g.startedAt = time.Now()
	slog.Info("Game loop started")

	// time.NewTicker(duração) = cria um "despertador" que toca a cada duração.
	// ticker.C é um CHANNEL (canal) = tubo de comunicação entre goroutines.
	// A cada intervalo, ticker envia um valor no channel C.
	ticker := time.NewTicker(g.interval)

	// defer ticker.Stop() = garante que o ticker é parado quando a função termina.
	// Parar é importante porque senão continua consumindo memória.
	defer ticker.Stop()

	// for { ... } = loop infinito.
	// Este loop roda enquanto o programa está ativo.
	// Sai quando ctx é cancelado (via select abaixo).
	for {
		// select = "aguarde um de vários canais ficar pronto e responda".
		// É como um switch, mas para o que chega primeiro de múltiplos canais.
		select {
		// case <-ctx.Done():
		// ctx.Done() retorna um canal que se activa quando o contexto é cancelado.
		// <-canal = "receba algo do canal" (bloqueia até algo chegar).
		// Se ctx foi cancelado, Done() retorna imediatamente.
		case <-ctx.Done():
			// Contexto foi cancelado (Ctrl+C ou Stop() foi chamado).
			// Muda running para false (não estamos mais rodando).
			g.running.Store(false)

			slog.Info("Game loop cancelled")

			// return = sai da função, o que finaliza a goroutine.
			return

		// case <-ticker.C:
		// ticker.C envia um sinal a cada 500ms (ou o que foi configurado).
		// Quando chega um tick, executamos 1 ciclo do jogo.
		case <-ticker.C:
			// Guarda a hora de ANTES de processar (para medir duração).
			start := time.Now()

			// Chama ProcessTick() do mundo = avança o jogo em 1 ciclo.
			// Usa o snapshot retornado — é o estado exato pós-tick, sem race condition.
			snap := g.world.ProcessTick()

			// Publica o snapshot para clientes WebSocket, se houver publisher.
			// Isso garante que o broadcast acontece APÓS todo o processamento do tick
			// (seguindo o ADR-005 — future managers como NPC vão antes desta linha).
			if g.publisher != nil {
				g.publisher.Publish(snap)
			}

			// Incrementa o contador de ticks de forma ATÔMICA.
			// Add(1) = adiciona 1 e retorna o novo valor (sem necessidade de lock).
			// atomic.Uint64 garante que isso é seguro mesmo com múltiplas goroutines lendo.
			g.tickCount.Add(1)

			// Calcula quanto tempo esse tick levou.
			// time.Since(t) = agora - t = duração desde 'start'.
			tickDuration := time.Since(start)

			// Trava o mutex para ESCRITA pois vamos modificar lastTickDuration.
			// Lock() = assegura que ninguém está lendo enquanto isso.
			g.mu.Lock()

			// Guarda a duração do último tick (para expor via Status()).
			g.lastTickDuration = tickDuration

			// Destrava o mutex.
			// Agora outras goroutines podem ler lastTickDuration.
			g.mu.Unlock()
		}
	}
}

// Stop para o game loop graciosamente.
// "Graciosamente" = com aviso prévio, permitindo que termine o que está fazendo.
func (g *GameLoop) Stop() {
	// CompareAndSwap(true, false) = só funciona se running está true
	// Se bem-sucedido, muda para false e retorna true
	// Se já estava false, retorna false
	if !g.running.CompareAndSwap(true, false) {
		// Se já não estava rodando, avisa e sai (não faz nada).
		slog.Warn("Game loop is not running")
		return
	}

	// Chama a função cancel() que foi guardada em Start().
	// Isso envia um sinal no context, ativando <-ctx.Done() no Start().
	// Quando ctx.Done() ativa, o loop faz return e finaliza.
	if g.cancel != nil {
		g.cancel()
	}

	slog.Info("Game loop stopped")
}

// Status retorna um snapshot do estado atual do loop.
// Uma "fotografia" que é thread-safe (segura para ler de múltiplas goroutines).
func (g *GameLoop) Status() Status {
	// g.mu.RLock() = trava para LEITURA compartilhada.
	// RLock = múltiplas goroutines podem ler ao mesmo tempo (porque ninguém está escrevendo).
	// Mas se houver um Lock (escrita) esperando, RLock fica bloqueado.
	g.mu.RLock()

	// Lê lastTickDuration de forma segura (protegido pelo lock acima).
	last := g.lastTickDuration

	// g.mu.RUnlock() = libera o lock de leitura.
	// Agora outra goroutine pode fazer Lock para escrever.
	g.mu.RUnlock()

	// Lê o horário de início do loop (não precisa de lock porque time.Time é simples).
	started := g.startedAt

	// Calcula uptime = quanto tempo o loop está rodando.
	uptime := time.Duration(0) // começa com 0
	if !started.IsZero() {     // se started não é data zero (1970-01-01)
		// uptime = agora - quando começou = tempo decorrido
		uptime = time.Since(started)
	}

	// Monta a resposta de status com TODAS as informações importantes.
	// Load() = lê um valor atômico de forma segura (sem lock).
	// Conversão int64(...) = converte uint64 para int64 (tipo esperado por Status).
	return Status{
		Running:          g.running.Load(),          // booleano: está rodando?
		Interval:         g.interval,                // tempo entre ticks
		TickCount:        int64(g.tickCount.Load()), // número total de ticks (com conversão de tipo)
		StartedAt:        started,                   // hora de início
		UpTime:           uptime,                    // quanto tempo rodando
		LastTickDuration: last,                      // quanto tempo o último tick levou
	}
}
