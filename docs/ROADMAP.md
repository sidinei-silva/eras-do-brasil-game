# Roadmap — Eras do Brasil

## Status Atual: Fase 0 — Heartbeat (reimplementando do zero)

| Fase | Nome | Objetivo | Status |
|------|------|----------|--------|
| **0** | Heartbeat | Servidor Go + game loop + WebSocket | Em andamento |
| **1** | Mundo Vivo | NPCs com Utility AI, ciclo dia/noite | Pendente |
| **2** | Observador | Cliente web read-only | Pendente |
| **3** | Jogador (MVP) | Criar, explorar, lutar, salvar | Pendente |
| **4** | Interacao | Dialogos, quests, crafting | Pendente |
| **5** | D20 Completo | 12 classes, grid tatico | Pendente |
| **6** | Multiplayer | Quests competitivas, temporadas | Pendente |

## Decisoes Tecnicas (ADRs)

| ADR | Decisao | Status |
|-----|---------|--------|
| [ADR-003](../decisions/ADR-003-estrategia-repositorios.md) | Repositorios separados (GDD / Game) | Ativo |
| [ADR-004](../decisions/ADR-004-pivot-mmorpg-servidor-go.md) | Pivot — Servidor Go + Cliente Web | Ativo |
| [ADR-005](../decisions/ADR-005-arquitetura-servidor-monolito-goroutines.md) | Monolito com goroutines | Ativo |
| [ADR-006](../decisions/ADR-006-persistencia-ram-first-sqlite.md) | RAM-first + SQLite | Ativo |

> ADRs de narrativa e mecanicas estao no repositorio [eras-do-brasil-gdd](https://github.com/sidinei-silva/eras-do-brasil-gdd).
