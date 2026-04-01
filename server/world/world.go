package world

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/sidinei-silva/eras-do-brasil-game/server/npc"
)

// Snapshot é a fotografia do estado do mundo em um momento.
// Os backticks com "json:" são tags que dizem qual nome usar quando converter para JSON.
type Snapshot struct {
	Tick      uint64    `json:"tick"`       // Número do ciclo do jogo
	GameTime  time.Time `json:"game_time"`  // Hora dentro do jogo
	Period    string    `json:"period"`     // Manhã/Tarde/Noite/Madrugada
	UpdatedAt time.Time `json:"updated_at"` // Quando esse snapshot foi criado
}

// World é o estado do mundo do jogo.
// O "mu" (mutex) protege os dados abaixo dele contra acesso simultâneo.
type World struct {
	// mu: MUTEX = fechadura que garante que só 1 goroutine acessa tick/gameTime por vez.
	// RWMutex = "Read-Write Mutex" = permite múltiplos leitores OU 1 escritor (não ambos).
	mu sync.RWMutex

	// tick: contador de quantas vezes ProcessTick foi chamado.
	tick uint64

	// gameTime: a hora atual dentro do jogo (aumenta a cada tick).
	gameTime time.Time

	WorldMap *WorldMap
	NPCs     *npc.Registry
}

// NewWorld cria uma nova instância do mundo.
// O símbolo * na frente de World significa que estamos retornando um PONTEIRO (endereço de memória)
// para um World, não uma cópia dele. Ponteiros permitem que mudanças afetam o original.
func NewWorld() *World {
	slog.Info("Creating new world")

	gameMap, err := LoadMataCosteira()
	if err != nil {
		panic(fmt.Sprintf("failed to load game map: %v", err))
	}
	slog.Info("Game map loaded", slog.String("region", gameMap.RegionName), slog.Int("zones", len(gameMap.Zones())))

	npcs, err := npc.LoadMataCosteira()
	if err != nil {
		panic(fmt.Sprintf("failed to load NPCs: %v", err))
	}
	slog.Info("NPCs loaded", slog.Int("count", npcs.Count()))

	return &World{
		gameTime: time.Date(1500, 1, 1, 6, 0, 0, 0, time.UTC),
		WorldMap: gameMap,
		NPCs:     npcs,
	}
}

// ProcessTick simula um ciclo do jogo: avança o tempo e incrementa contador.
// O (w *World) significa: "este método tem receiver w, que é um PONTEIRO para World".
// receiver = quem chama o método. Se é *, significal que modifica o original, não uma cópia.
func (w *World) ProcessTick() Snapshot {
	// w.mu.Lock() = TRAVA a fechadura para ESCRITA exclusiva.
	// Isso significa: "nenhuma outra goroutine pode ler ou escrever até eu chamar Unlock".
	// Se outra goroutine estiver dentro de Lock/Unlock, esta goroutine ESPERA aqui até ela sair.
	// REGRA IMPORTANTE: Lock() DEVE ser combinado com Unlock() (não RUnlock()).
	w.mu.Lock()

	// defer = "execute isto quando a função terminar, não importa o que aconteça".
	// defer w.mu.Unlock() = garante que o lock é solto mesmo se houver erro ou panic.
	// Isso é MUITO importante para evitar deadlock (travar para sempre).
	defer w.mu.Unlock()

	// Incrementa o contador de ticks do jogo.
	// w.tick++ = "pegue o valor de w.tick, some 1, e guarde de volta".
	w.tick++
	w.gameTime = w.gameTime.Add(1 * time.Minute)
	period := getPeriod(w.gameTime.Hour())

	// Processa todos os NPCs neste tick.
	// Os NPCs estão protegidos pelo mesmo mutex do World — acesso seguro.
	w.NPCs.ProcessTick(period)

	slog.Info("Processing tick", slog.Uint64("tick", w.tick))

	return Snapshot{
		Tick:      w.tick,
		GameTime:  w.gameTime,
		Period:    period,
		UpdatedAt: time.Now().UTC(),
	}
}

func (w *World) Snapshot() Snapshot {
	// w.mu.RLock() = TRAVA a fechadura para LEITURA.
	// Permite que múltiplas goroutines leiam ao mesmo tempo, mas bloqueia se alguém estiver escrevendo.
	w.mu.RLock()

	// defer w.mu.RUnlock() = garante que o lock de leitura é solto quando a função terminar.
	defer w.mu.RUnlock()

	// Retorna uma fotografia do estado ATUAL do mundo.
	return Snapshot{
		Tick:      w.tick,
		GameTime:  w.gameTime,
		Period:    getPeriod(w.gameTime.Hour()),
		UpdatedAt: time.Now().UTC(),
	}
}

// Períodos do dia — nomes canônicos em português, alinhados com o GDD.
// Usar estas constantes em vez de strings literais evita typos e facilita busca no código.
const (
	PeriodoManha     = "Manha"     // 06:00 – 11:59
	PeriodoTarde     = "Tarde"     // 12:00 – 17:59
	PeriodoNoite     = "Noite"     // 18:00 – 21:59
	PeriodoMadrugada = "Madrugada" // 22:00 – 05:59
)

// getPeriod retorna qual período do dia é, baseado na hora (0-23).
// Não tem receiver (não está acoplado a World), então é uma função simples.
func getPeriod(h int) string {
	// switch = escolhe uma ação baseado no valor de h.
	// case = "se o valor for isto, faça aquilo".
	switch {
	case h >= 6 && h < 12:
		return PeriodoManha
	case h >= 12 && h < 18:
		return PeriodoTarde
	case h >= 18 && h < 22:
		return PeriodoNoite
	default: // 22-05
		return PeriodoMadrugada
	}
}
