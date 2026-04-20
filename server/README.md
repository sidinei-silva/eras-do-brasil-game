# 🖥️ Server — Eras do Brasil

> Servidor Go do MUD Moderno (MMORPG server-authoritative).

## Stack

- **Go 1.22+** — Goroutines para game loop, IA de NPCs e conexões simultâneas
- **WebSocket** (gorilla/websocket) — Comunicação bidirecional em tempo real
- **SQLite** (modernc.org/sqlite) — Persistência RAM-First com snapshots async

## Arquitetura

- [ADR-005 — Monolito + Game Loop Sequencial](../decisions/ADR-005-arquitetura-servidor-monolito-goroutines.md)
- [ADR-006 — Persistência RAM-First + SQLite](../decisions/ADR-006-persistencia-ram-first-sqlite.md)
- [ADR-007 — Estrutura de Pacotes (flat por domínio)](../decisions/ADR-007-estrutura-pacotes-server.md)
- [Documentação técnica de arquitetura](../docs/architecture/00-index.md)

## Estrutura Planejada

```
server/
├── main.go              # entry point, wiring, graceful shutdown
├── engine/              # GameLoop, CommandQueue
├── state/               # GameState — estado mutável central
├── command/             # tipos de comando + roteadores
├── world/               # WorldManager + GameTime + Block
├── npc/                 # NPCManager + NPC + Agenda + Need
├── net/                 # HTTP + WebSocket
│   └── socket/
│       ├── player/
│       └── admin/
├── data/                # loaders JSON + templates estáticos
└── persist/             # SQLite (Fase 1+)
```

### Regras de dependência

1. `main.go` é o único que importa todos os pacotes
2. Managers (`world`, `npc`, ...) **não** importam outros managers —
   trocam dados via `state.GameState`
3. `net/` e `persist/` não conhecem lógica de jogo
4. Pacote nasce quando tem habitante (sem pacotes vazios "por precaução")

## Comandos Admin (Dev/Teste)

- **Propósito:** acelerar debug e validação sem depender da interface final.
- **Arquitetura:** `AdminCommandManager` roda em goroutine dedicada e integra com EventBus.
- **Fluxo:** entrada de comando -> validação/permissão -> despacho para manager de domínio -> resposta -> auditoria.
- **Entrada por fase:**
  - Fase 0-2: console local
  - Fase 3-5: endpoint interno administrativo
  - Fase 6+: remoto opcional (RCON-like), se necessário
- **Segurança:** comandos destrutivos somente em dev/homolog e com trilha de auditoria.

## Como rodar

```bash
cd server
go run .
# Acesse http://localhost:8080
```

## Fase Atual

**Fase 0: Heartbeat** — Game loop + WebSocket + cliente mínimo.
