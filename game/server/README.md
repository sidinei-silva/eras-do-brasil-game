# 🖥️ Server — Eras do Brasil

> Servidor Go do MUD Moderno (MMORPG server-authoritative).

## Stack

- **Go 1.22+** — Goroutines para tick loop, IA de NPCs e conexões simultâneas
- **WebSocket** (gorilla/websocket) — Comunicação bidirecional em tempo real
- **SQLite** (modernc.org/sqlite) — Persistência RAM-First com snapshots async

## Arquitetura

Definida nos ADRs:

- [ADR-005 — Monolito + Goroutines + EventBus](../../decisions/ADR-005-arquitetura-servidor-monolito-goroutines.md)
- [ADR-006 — Persistência RAM-First + SQLite](../../decisions/ADR-006-persistencia-ram-first-sqlite.md)

## Estrutura Planejada

```
server/
├── main.go           # Entry point, HTTP server, WebSocket upgrade
├── go.mod
├── engine/
│   ├── eventbus.go   # Pub/Sub com Go channels
│   └── tick.go       # TickEngine com time.Ticker
├── world/            # WorldManager — blocos, grafo, dia/noite
├── combat/           # CombatManager — D20, turnos, loot
├── economy/          # EconomyManager — crafting, comércio
├── player/           # Conexão, readPump/writePump, estado do jogador
├── narrative/        # StoryManager, NPCManager, diálogos
├── admin/            # AdminCommandManager (console dev, comandos admin)
├── persist/          # PersistManager — SQLite snapshots
└── data/             # JSON templates (itens, classes, inimigos)
```

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
cd game/server
go run main.go
# Acesse http://localhost:8080
```

## Fase Atual

**Fase 0: Heartbeat** — Tick loop + WebSocket + cliente mínimo.
