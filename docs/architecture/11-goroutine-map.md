# 11 — Mapa de Goroutines Ativas

> **Quando consultar:** Quando quiser ter certeza de que não está perdendo goroutines (leak), quando precisar entender quantas goroutines existem ao mesmo tempo, ou quando uma goroutine não estiver terminando no shutdown.

## Diagrama — Todas as goroutines em runtime

```mermaid
graph TD
    subgraph "Goroutine principal (main)"
        MAIN["main()<br/><i>Bloqueia em ctx.Done()</i><br/><i>Depois espera wg.Wait()</i>"]
    end

    subgraph "Goroutines de longa duração (vivem o servidor inteiro)"
        GL["GameLoop<br/><i>1 goroutine</i><br/>Criada: main → go gameLoop.Start(ctx)"]
        EB["EventBus<br/><i>1 goroutine</i><br/>Criada: main → go eventBus.Run(ctx)"]
        PM["PersistManager<br/><i>1 goroutine</i><br/>Criada: main → go persist.Run(ctx)"]
        HUB["Hub<br/><i>1 goroutine</i><br/>Criada: main → go hub.Run(ctx)"]
    end

    subgraph "Goroutines por jogador (vivem enquanto o jogador está conectado)"
        RP1["readPump jogador A<br/><i>1 goroutine</i>"]
        WP1["writePump jogador A<br/><i>1 goroutine</i>"]
        RP2["readPump jogador B<br/><i>1 goroutine</i>"]
        WP2["writePump jogador B<br/><i>1 goroutine</i>"]
        RPN["readPump jogador N<br/><i>...</i>"]
        WPN["writePump jogador N<br/><i>...</i>"]
    end

    MAIN -->|"go"| GL
    MAIN -->|"go"| EB
    MAIN -->|"go"| PM
    MAIN -->|"go"| HUB

    HUB -->|"nova conexão → go"| RP1
    HUB -->|"nova conexão → go"| WP1
    HUB -->|"nova conexão → go"| RP2
    HUB -->|"nova conexão → go"| WP2
```

## Contagem de goroutines

Com N jogadores conectados, o servidor tem exatamente:

| Goroutine | Quantidade | Quem cria | Quem mata |
|-----------|-----------|-----------|-----------|
| main | 1 | runtime | SIGTERM → retorna |
| GameLoop | 1 | main | ctx cancelado |
| EventBus | 1 | main | ctx cancelado |
| PersistManager | 1 | main | ctx cancelado |
| Hub | 1 | main | ctx cancelado |
| readPump | N (1 por jogador) | Hub (na conexão) | Erro de leitura ou ctx cancelado |
| writePump | N (1 por jogador) | Hub (na conexão) | ctx cancelado ou readPump morreu |

**Total: 5 + 2N goroutines.**

Com 10 jogadores: 25 goroutines. Com 100 jogadores: 205 goroutines. Goroutines em Go são leves (~8KB de stack cada), então 205 goroutines consomem ~1.6MB. Irrelevante.

## Diagrama — Ciclo de vida e cancelamento

```mermaid
sequenceDiagram
    participant MAIN as main()
    participant CTX as context (raiz)
    participant GL as GameLoop
    participant HUB as Hub
    participant PM as PersistManager
    participant RP as readPumps
    participant WP as writePumps

    Note over MAIN: Servidor rodando normalmente

    MAIN->>CTX: signal.NotifyContext(SIGTERM)
    Note over CTX: Ctrl+C → ctx cancelado

    CTX->>GL: ctx.Done() → GameLoop para o ticker
    CTX->>HUB: ctx.Done() → Hub para de aceitar conexões
    CTX->>PM: ctx.Done() → PersistManager faz flush final

    HUB->>RP: Fecha conexões → readPumps terminam
    RP->>WP: readPump morreu → writePump detecta e termina

    GL->>MAIN: wg.Done()
    HUB->>MAIN: wg.Done()
    PM->>MAIN: wg.Done()

    MAIN->>MAIN: wg.Wait() desbloqueia → os.Exit(0)
```

## Explicação — Como evitar goroutine leak

Goroutine leak é quando uma goroutine fica rodando pra sempre sem ninguém ouvindo, consumindo memória e CPU. As causas mais comuns:

**Causa 1: Canal sem leitor.** Se uma goroutine escreve num channel e ninguém lê, ela bloqueia pra sempre. Solução: channels buffered ou `select` com `default` ou `ctx.Done()`.

**Causa 2: Loop infinito sem condição de saída.** Se o readPump tem `for { conn.Read() }` sem checar o ctx, e o WebSocket nunca retorna erro, a goroutine nunca termina. Solução: passar ctx para `conn.Read(ctx)` — quando ctx cancela, Read retorna erro.

**Causa 3: defer que nunca executa.** Se uma goroutine entra num loop antes do defer ser registrado. Solução: defer logo no início da função, antes de qualquer loop.

**Regra prática:** Toda goroutine de longa duração deve ter uma das seguintes condições de saída no seu loop:
- `case <-ctx.Done():` no select
- Erro retornado por uma operação bloqueante (Read, Write)
- Channel fechado que o `range` detecta

## Monitorando goroutines

Em desenvolvimento, adicione ao endpoint `/admin/status`:

```json
{
  "goroutines": 25,
  "players_online": 10,
  "tick_count": 14832,
  "last_tick_duration_ms": 3
}
```

O valor de `runtime.NumGoroutine()` do Go retorna o total. Se esse número só cresce e nunca diminui (mesmo quando jogadores desconectam), há um leak.

## Quando a goroutine certa é NENHUMA

Nem tudo precisa de goroutine. Regras:

- **Lógica de jogo (managers)**: NÃO usa goroutine. Roda sequencialmente dentro do tick do game loop.
- **I/O de longa duração (WebSocket, persistência)**: USA goroutine. Fica bloqueada esperando I/O a maior parte do tempo.
- **Cálculo curto dentro do tick (rolar D20, mover NPC)**: NÃO usa goroutine. É rápido o suficiente pra rodar no tick.
- **Cálculo pesado que excede tick budget (futuro: pathfinding de 500 NPCs)**: USA goroutine worker pool DENTRO do manager, não fora dele.
