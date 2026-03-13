# Eras do Brasil — Game

> MUD Moderno (MMORPG server-authoritative) — Servidor Go + Cliente Web

## Sobre

Repositório de **código** do projeto Eras do Brasil. Aqui ficam o servidor Go, o cliente web e toda a documentação técnica de desenvolvimento.

O **GDD (Game Design Document)** e o lore estão no repositório [eras-do-brasil-gdd](https://github.com/sidinei-silva/eras-do-brasil-gdd).

## Stack

- **Go 1.22+** — Goroutines para tick loop, IA de NPCs e conexões simultâneas
- **WebSocket** (gorilla/websocket) — Comunicação bidirecional em tempo real
- **SQLite** (modernc.org/sqlite) — Persistência RAM-First com snapshots async
- **Cliente Web** — HTML/CSS/JS puro via WebSocket

## Estrutura

```
├── game/
│   ├── server/       # Servidor Go (engine, world, combat, economy, etc.)
│   └── textClient/   # Cliente web HTML/CSS/JS
├── decisions/        # ADRs (Architecture Decision Records)
├── docs/             # Guias e specs de produto
├── historico/        # Sessões de trabalho anteriores
├── ROADMAP.md        # Fases de desenvolvimento
├── backlog.md        # Tarefas detalhadas
└── project-status.md # Status atual do projeto
```

## ADRs

- [ADR-003 — Estratégia de Repositórios](decisions/ADR-003-estrategia-repositorios.md)
- [ADR-004 — Pivot MMORPG Servidor Go](decisions/ADR-004-pivot-mmorpg-servidor-go.md)
- [ADR-005 — Arquitetura Monolito + Goroutines + EventBus](decisions/ADR-005-arquitetura-servidor-monolito-goroutines.md)
- [ADR-006 — Persistência RAM-First + SQLite](decisions/ADR-006-persistencia-ram-first-sqlite.md)

## Links

- **GDD & Lore:** [eras-do-brasil-gdd](https://github.com/sidinei-silva/eras-do-brasil-gdd)
- **Legado (Unity-era):** [ErasDoBrasil-Historico-Legado](https://github.com/sidinei-silva/ErasDoBrasil-Historico-Legado)
- **Site:** [erasdobrasil.com.br](https://erasdobrasil.com.br)
