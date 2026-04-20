# 🖥️ Server — Eras do Brasil

> Servidor Go do MUD Moderno (MMORPG server-authoritative).

## Stack

- **Go 1.22+** — Goroutines para game loop, IA de NPCs e conexões simultâneas
- **WebSocket** (gorilla/websocket) — Comunicação bidirecional em tempo real
- **SQLite** (modernc.org/sqlite) — Persistência RAM-First com snapshots async (Fase 1+)

## Arquitetura

- [ADR-005 — Monolito + Game Loop Sequencial](../decisions/ADR-005-arquitetura-servidor-monolito-goroutines.md)
- [ADR-006 — Persistência RAM-First + SQLite](../decisions/ADR-006-persistencia-ram-first-sqlite.md)
- [ADR-007 — Estrutura de Pacotes (flat por domínio, managers donos do próprio estado)](../decisions/ADR-007-estrutura-pacotes-server.md)
- [Documentação técnica de arquitetura](../docs/architecture/00-index.md)

## Estrutura

```
server/
├── main.go              # entry point, wiring, graceful shutdown
├── engine/              # GameLoop
├── command/             # PlayerCommand, CommandQueue, routers
├── world/               # Manager (dono do GameTime; futuro: Block, Climate)
├── npc/                 # Manager (dono dos NPCs; Need, Schedule, Behavior)
├── snapshot/            # Snapshot — agregador transversal (cópia pontual)
├── net/                 # HTTP + WebSocket
│   ├── api/
│   └── socket/
│       ├── player/
│       └── admin/
├── data/                # JSON templates (NPCs, futuramente items/quests)
└── persist/             # SQLite (Fase 1+)
```

Pacotes futuros (`mob/`, `combat/`, `story/`, `economy/`) nascem na raiz
quando forem implementados — cada um dono do próprio estado, expondo
`ProcessTick` e getters para o snapshot.

### Regras de dependência

1. **`main.go` é o único que importa todos os managers** — faz wiring
   e injeta dependências no wire-up.
2. **Managers são donos do próprio estado.** Não existe god object central.
   `world.Manager` é dono do `GameTime`. `npc.Manager` é dono dos NPCs.
   Cada manager futuro vai ser dono do seu domínio.
3. **Managers não importam comportamento de outros managers.** Podem
   importar **tipos universais** (ex: `npc` importa `world.GameTime` pra
   usar na assinatura do `ProcessTick`) — isso é acoplamento a dado,
   não a lógica.
4. **`snapshot/` é o único pacote que importa todos os managers.** Ele
   existe pra produzir uma cópia imutável do estado a cada tick, que vai
   pro Admin Hub e (futuramente) pro Persist. Nenhum manager importa
   `snapshot`.
5. **`net/` e `persist/` não conhecem lógica de jogo.** Recebem
   `PlayerCommand` (via `command.Queue`) e entregam `Snapshot` pra
   serializar.
6. **Pacote nasce quando tem habitante.** Não cria pacote vazio "por
   precaução" — criar quando tiver código real pra pôr dentro.

### Fluxo do tick

GameLoop (a cada 1s):

1. command.ProcessPlayerCommands(tickCount, pendingCommands)
2. worldManager.ProcessTick()→ avança GameTime
3. npcManager.ProcessTick(worldManager.GameTime())→ NPCs veem tempo atual
4. snap := snapshot.Build(tickCount, worldManager, npcManager)
5. adminHub.Publish(snap)

## Comandos Admin (Dev/Teste)

- **Propósito:** acelerar debug e validação sem depender da interface final.
- **Arquitetura:** `AdminRouter` recebe comandos do WebSocket admin e
  decide: comandos de mutação vão pra `CommandQueue` (processados no tick),
  comandos de consulta são respondidos out-of-band (OOB) direto ao cliente.
- **Fluxo:** entrada → validação → fila (mutação) OU resposta direta (consulta) → auditoria.
- **Entrada por fase:**
  - Fase 0-2: cliente web admin (já implementado)
  - Fase 3-5: endpoint administrativo autenticado
  - Fase 6+: remoto opcional (RCON-like), se necessário
- **Segurança:** comandos destrutivos somente em dev/homolog e com trilha de auditoria.

## Como rodar

```bash
cd server
go run .
# Acesse http://localhost:8080
```

## Fase Atual

**Fase 0: Heartbeat** — Game loop + WebSocket + cliente admin + NPCs
com rotina por período do dia + necessidades (fome/energia).
