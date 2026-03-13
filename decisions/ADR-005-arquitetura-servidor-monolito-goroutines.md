# ADR-005: Arquitetura do Servidor — Monolito com Goroutines e EventBus

## Status

Aceito

## Contexto

Com o pivot para servidor Go (ADR-004), precisamos definir como os "motores" do jogo (tick, mundo, NPCs, combate, narrativa, economia) se comunicam entre si.

Duas abordagens possíveis:

1. **N processos (microserviços)** — cada motor é um serviço separado, comunicação via gRPC/TCP/Redis
2. **1 processo, N goroutines** — cada motor é uma goroutine, comunicação via channels Go

O developer é solo, tem 6h/semana, e está aprendendo Go. A complexidade operacional precisa ser mínima.

## Decisão

**1 processo, N goroutines com EventBus baseado em channels Go.**

### Arquitetura de Goroutines

```
main.go (startup + wiring + graceful shutdown)
  │
  ├── TickEngine        goroutine — time.Ticker dispara ticks
  ├── WorldManager      goroutine — processa estado do mundo
  ├── NPCManager        goroutine — rotinas, IA, fofoca
  ├── CombatManager     goroutine — combates ativos
  ├── StoryManager      goroutine — quests, eventos, ruptura
  ├── EconomyManager    goroutine — preços, comércio
  ├── PersistManager    goroutine — snapshots async (ver ADR-006)
  ├── EventBus          goroutine — roteamento pub/sub
  │
  └── Per-connection (1 par por jogador):
      ├── readPump      goroutine — lê WebSocket → comandos
      └── writePump     goroutine — envia estado → WebSocket
```

### EventBus (Pub/Sub com Channels)

- `Publish(topic string, event Event)` — broadcast para todos os subscribers do tópico
- `Subscribe(topic string) <-chan Event` — retorna channel de leitura
- `Unsubscribe(topic string, ch <-chan Event)` — remove subscriber
- Thread-safe via `sync.RWMutex` no mapa de subscribers
- Cada manager faz `for event := range channel` — bloqueia sem polling, zero CPU idle

### Tópicos de Eventos (exemplos iniciais)

| Tópico | Produtor | Consumidores |
|--------|----------|-------------|
| `tick` | TickEngine | World, NPC, Combat, Story, Economy |
| `player.entered` | Session | World, NPC, Story |
| `player.moved` | World | NPC, Story |
| `npc.spoke` | NPC | Story, Session (nearby players) |
| `combat.started` | Combat | NPC, Story, Session |
| `combat.ended` | Combat | World, NPC, Persist |
| `item.traded` | Economy | Persist, Story |
| `quest.completed` | Story | NPC, Economy, Persist |

### Fluxo de Exemplo

```
TickEngine dispara "tick"
  → WorldManager recebe: atualiza dia/noite, gera "period_changed"
  → NPCManager recebe "tick": executa rotinas, gera "npc_moved", "npc_spoke"
  → StoryManager recebe "npc_spoke": avalia trigger de quest
  → PersistManager recebe dirty entities via channel dedicado
  → PlayerConn recebe eventos relevantes → envia JSON via WebSocket
```

### Estrutura de Módulos

```
server/
├── main.go              ← startup, wiring, graceful shutdown
├── engine/
│   ├── eventbus.go      ← pub/sub de eventos
│   └── tick.go          ← TickEngine (time.Ticker)
├── world/
│   ├── mundo.go         ← struct Mundo, ProcessarTick()
│   ├── bloco.go         ← struct Bloco (nó do mapa)
│   └── npc.go           ← struct NPC, rotinas, Utility AI
├── combat/
│   ├── manager.go       ← CombatManager
│   └── d20.go           ← motor de rolagem
├── economy/
│   ├── inventario.go
│   └── comercio.go
├── player/
│   ├── player.go        ← struct Player
│   └── session.go       ← WebSocket connection handler
├── narrative/
│   ├── story.go         ← StoryManager
│   └── quest.go
├── persist/
│   └── persist.go       ← snapshot goroutine (ver ADR-006)
├── data/
│   └── *.json           ← dados estáticos (NPCs, items, mapa)
└── go.mod
```

### Quando Considerar Separação de Processos

Somente se **todas** estas condições forem verdadeiras:

1. >500 jogadores simultâneos
2. Profiling (`pprof`) mostra CPU saturada em um manager específico
3. O manager gargalo pode rodar independente sem acesso à memória compartilhada

Até lá, goroutines + channels é a solução correta.

## Consequências

**Positivas:**

- Memória compartilhada — sem serialização/deserialização entre motores
- Debugging simples — 1 processo, 1 stack trace, `go tool pprof`
- Zero infraestrutura — sem Docker, sem message broker, sem service discovery
- Channels nativos do Go — concorrência sem reinventar a roda
- EventBus desacopla managers — adicionar um novo motor = subscribe nos tópicos relevantes

**Negativas:**

- Se um manager fizer `panic` não recuperado, derruba o processo inteiro (mitigado: `recover()` em cada goroutine)
- Escalabilidade vertical tem limite físico (mitigado: irrelevante para solo dev com <100 jogadores)
- Toda goroutine compartilha o mesmo espaço de memória — precisa disciplina com ponteiros (mitigado: EventBus copia eventos, managers operam em seus próprios dados)

## Alternativas Consideradas

- **Microserviços (1 processo por motor):** Rejeitado — complexidade operacional desproporcional para solo dev. Serialização JSON/gRPC entre processos adiciona latência e boilerplate. Estado distribuído exige Redis/etcd. Debugging distribuído é ordens de magnitude mais difícil.
- **Actor model (e.g. Proto.Actor):** Rejeitado — adiciona dependência externa e abstração desnecessária quando channels Go nativos resolvem o mesmo problema com menos magia.
- **Single goroutine (game loop síncrono):** Rejeitado — não aproveita multi-core, NPC AI pesada bloquearia o tick.

## Referências

- ADR-004: Pivot para MMORPG — define stack Go + WebSocket
- ADR-006: Estratégia de Persistência — define o PersistManager mencionado aqui
- GDD Cap 8 §8.12: Arquitetura do Motor de Mundo
