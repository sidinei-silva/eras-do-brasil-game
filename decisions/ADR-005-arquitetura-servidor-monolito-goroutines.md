# ADR-005: Arquitetura do Servidor — Monolito com Game Loop Sequencial

## Status

Aceito (revisado em 2026-03-15 — corrige arquitetura de tick paralelo para game loop sequencial)

## Contexto

Com o pivot para servidor Go (ADR-004), precisamos definir como os "motores" do jogo (mundo, NPCs, combate, narrativa, economia) se comunicam entre si.

Três abordagens possíveis:

1. **N processos (microserviços)** — cada motor é um serviço separado, comunicação via gRPC/TCP/Redis
2. **N goroutines paralelas com EventBus** — cada motor é uma goroutine, tick broadcast via pub/sub
3. **Game loop sequencial + goroutines de I/O** — lógica de jogo em uma goroutine, I/O em paralelo

O developer é solo, tem 6h/semana, e está aprendendo Go. A complexidade operacional precisa ser mínima.

## Decisão

**1 processo. Game loop sequencial para lógica de jogo. Goroutines apenas para I/O.**

Este é o padrão da indústria para game servers de simulação de mundo (World of Warcraft, RuneScape, Minecraft). A lógica de jogo executa sequencialmente dentro de cada tick para garantir consistência de estado — NPCManager vê o estado que WorldManager acabou de atualizar no mesmo tick.

### Arquitetura

```
main.go (startup + wiring + graceful shutdown)
  │
  ├── GameLoop            goroutine PRINCIPAL — sequencial, processa todos os managers por tick
  │     │
  │     ├── world.ProcessTick(gameState)     ← atualiza dia/noite, clima, blocos
  │     ├── npc.ProcessTick(gameState)       ← rotinas, IA, fofoca (vê estado atualizado do world)
  │     ├── combat.ProcessTick(gameState)    ← resolve turnos de combates ativos
  │     ├── story.ProcessTick(gameState)     ← avalia triggers, quests, arcos
  │     ├── economy.ProcessTick(gameState)   ← preços, comércio pendente
  │     └── persist.QueueDirty(gameState)    ← enfileira entidades sujas para goroutine de I/O
  │
  ├── PersistManager      goroutine de I/O — snapshots async (ver ADR-006)
  ├── AdminManager        goroutine de I/O — comandos administrativos
  ├── EventBus            goroutine — notificações assíncronas (clientes, logging, métricas)
  │
  └── Per-connection (1 par por jogador):
      ├── readPump        goroutine — lê WebSocket → fila de comandos para o GameLoop
      └── writePump       goroutine — envia estado → WebSocket
```

### Por que sequencial e não paralelo

O problema com managers paralelos recebendo tick via EventBus:

```
❌ Managers paralelos (versão anterior deste ADR):
  TickEngine dispara "tick" via broadcast
    → WorldManager recebe: atualiza dia/noite
    → NPCManager recebe "tick" SIMULTANEAMENTE: executa rotinas
    ⚠️ NPCManager pode NÃO ver a mudança de período que WorldManager fez
    ⚠️ Ordem de processamento indeterminada = bugs intermitentes
```

```
✅ Game loop sequencial (padrão da indústria):
  GameLoop a cada tick:
    world.ProcessTick()   → dia vira noite
    npc.ProcessTick()     → NPCs veem que é noite, mudam comportamento
    combat.ProcessTick()  → combates usam estado correto
    story.ProcessTick()   → narrativa reage a tudo que aconteceu neste tick
    ✓ Ordem garantida, estado consistente, zero race conditions
```

### Tick Budget

Com tick de 500ms, o budget total para ProcessTick de todos os managers é 500ms. Se o processamento exceder o budget:

- **Fase 0-3** (~100 NPCs): irrelevante — processamento leva <10ms
- **Fase 4-5** (~500 NPCs): monitorar via `pprof`, otimizar hot paths
- **Se exceder**: paralelizar DENTRO de um manager (ex: dividir NPCs em batches processados por worker goroutines), mas manter a ORDEM entre managers sequencial

### EventBus (papel revisado)

O EventBus continua existindo, mas **não coordena lógica de tick**. Seu papel é:

- Notificar clientes WebSocket de eventos relevantes (NPC falou, combate começou)
- Logging e métricas
- Eventos assíncronos que não precisam de ordem (conquistas, notificações)

Interface:

- `Publish(topic string, event Event)` — broadcast para subscribers
- `Subscribe(topic string) <-chan Event` — retorna channel de leitura
- Subscrições acontecem **apenas no startup** (invariante: sem Subscribe/Unsubscribe em runtime)

### Tópicos de Eventos (assíncronos — notificação, não coordenação)

| Tópico | Produtor | Consumidores |
|--------|----------|-------------|
| `player.entered` | GameLoop | Session (nearby players) |
| `npc.spoke` | GameLoop | Session (nearby players) |
| `combat.started` | GameLoop | Session (participants) |
| `combat.ended` | GameLoop | Session (participants), Persist |
| `item.traded` | GameLoop | Persist |
| `quest.completed` | GameLoop | Session (player) |

### Fluxo de Exemplo (revisado)

```
time.Ticker dispara → GameLoop acorda
  1. world.ProcessTick(state)     → dia vira noite, gera evento "period_changed"
  2. npc.ProcessTick(state)       → NPCs veem noite, vão dormir, gera "npc_moved"
  3. combat.ProcessTick(state)    → resolve turno do combate ativo
  4. story.ProcessTick(state)     → avalia se threshold de quest foi atingido
  5. economy.ProcessTick(state)   → processa trades pendentes
  6. persist.QueueDirty(state)    → enfileira entidades modificadas

  EventBus.Publish("npc_moved", ...) → writePump envia para clientes nearby
  EventBus.Publish("period_changed", ...) → writePump atualiza HUD de todos
```

### Fila de Comandos do Jogador

Comandos recebidos via WebSocket (readPump) vão para uma fila thread-safe. O GameLoop consome essa fila **no início de cada tick**, antes de processar os managers:

```
GameLoop a cada tick:
  1. Processar fila de comandos do jogador (mover, atacar, usar item)
  2. world.ProcessTick()
  3. npc.ProcessTick()
  4. ...
```

Isso garante que ações do jogador são processadas no contexto correto do tick.

### Estrutura de Módulos

> A organização concreta dos pacotes está definida no [ADR-007](ADR-007-estrutura-pacotes-server.md).
> Resumo: layout flat por domínio na raiz de `server/`. Cada manager é um
> pacote (`world/`, `npc/`, `combat/`, ...). Estado mutável central em
> `state/`. I/O em `net/` e `persist/`. Sem `internal/`, sem `cmd/`.

A relação entre os pacotes segue duas regras invioláveis:

1. **Managers não importam outros managers** — comunicação via
   `state.GameState` que o GameLoop passa em cada `ProcessTick()`.
2. **Pacotes de I/O não conhecem lógica de jogo** — recebem comandos
   via `command.Queue`, devolvem resultados via channels ou snapshots.

### Quando Considerar Paralelismo no Game Loop

Somente se **todas** estas condições forem verdadeiras:

1. `pprof` mostra que um `ProcessTick()` específico excede 50% do tick budget
2. O manager pode ser paralelizado internamente (ex: dividir NPCs em batches)
3. O paralelismo é DENTRO do manager, não ENTRE managers (a ordem sequencial entre managers é inviolável)

## Consequências

**Positivas:**

- Estado consistente dentro de cada tick — zero race conditions na lógica de jogo
- Debugging trivial — stack trace linear, comportamento determinístico
- Zero infraestrutura — sem Docker, sem message broker, sem service discovery
- Padrão da indústria — mesma arquitetura de WoW, RuneScape, Minecraft
- Goroutines de I/O aproveitam multi-core onde faz diferença (WebSocket, persistência)

**Negativas:**

- Lógica de jogo roda em 1 core (mitigado: irrelevante para <1000 NPCs com tick de 500ms)
- Se um manager fizer `panic` não recuperado, derruba o processo inteiro (mitigado: `recover()` no GameLoop)
- Tick budget é compartilhado entre todos os managers (mitigado: monitorar com `pprof`, otimizar antes de paralelizar)

## Alternativas Consideradas

- **Microserviços (1 processo por motor):** Rejeitado — complexidade operacional desproporcional para solo dev. Serialização JSON/gRPC entre processos adiciona latência e boilerplate. Estado distribuído exige Redis/etcd. Debugging distribuído é ordens de magnitude mais difícil.
- **Actor model (e.g. Proto.Actor):** Rejeitado — adiciona dependência externa e abstração desnecessária quando channels Go nativos resolvem o mesmo problema com menos magia.
- **N goroutines paralelas com EventBus para tick (versão anterior deste ADR):** Rejeitado — managers paralelos recebendo tick via broadcast não têm ordem determinística. NPCManager pode processar antes de WorldManager atualizar dia/noite. Gera race conditions sutis entre managers que compartilham estado. Anti-pattern para game servers — a indústria usa game loop sequencial.

## Referências

- ADR-004: Pivot para MMORPG — define stack Go + WebSocket
- ADR-006: Estratégia de Persistência — define o PersistManager mencionado aqui
- GDD Cap 8 §8.12: Arquitetura do Motor de Mundo
- Padrão "Game Loop": [Game Programming Patterns — Game Loop](https://gameprogrammingpatterns.com/game-loop.html)
- ADR-007:[Estrutura de Pacotes do Servidor](./ADR-007-estrutura-pacotes-server.md)
