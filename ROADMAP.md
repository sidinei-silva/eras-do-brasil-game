# Roadmap — Eras do Brasil

> Para retomar o desenvolvimento, use o [Guia de Retomada](docs/GUIA_RETOMADA.md).

---

## Status Atual: Fase 0 — Heartbeat

---

## Fases de Desenvolvimento

| Fase | Nome | Objetivo | Status |
|------|------|----------|--------|
| **0** | Heartbeat | Servidor Go com game loop + WebSocket + cliente mínimo | 🟡 Próxima |
| **1** | Living World | Blocos, NPCs com Utility AI, ciclo dia/noite (relógio do jogo) | 🔲 Pendente |
| **2** | Observer | Cliente web observa o mundo em tempo real (read-only) | 🔲 Pendente |
| **3** | Player ≈ MVP | Criação de personagem, D20, combate, inventário, save | 🔲 Pendente |
| **4** | Interaction | Diálogos, facções, crafting, temporadas | 🔲 Pendente |
| **5** | D20 Full | 12 classes, 3 Origens, grid tático, Ato 1 completo | 🔲 Pendente |
| **6** | Multiplayer | Penalidade de morte, quests competitivas, inimigos evolutivos, temporadas | 🔲 Pendente |

```
Fase 0 ──► Fase 1 ──► Fase 2 ──► Fase 3 ──► Fase 4 ──► Fase 5 ──► Fase 6
```

---

## Decisões Técnicas (ADRs)

| ADR | Decisão | Status |
|-----|---------|--------|
| [ADR-003](decisions/ADR-003-estrategia-repositorios.md) | Estratégia de repositórios (GDD separado do Game) | Ativo |
| [ADR-004](decisions/ADR-004-pivot-mmorpg-servidor-go.md) | Pivot — Servidor Go + Cliente Web (MMORPG) | **Ativo** |
| [ADR-005](decisions/ADR-005-arquitetura-servidor-monolito-goroutines.md) | Arquitetura do servidor: monólito com goroutines | Ativo |
| [ADR-006](decisions/ADR-006-persistencia-ram-first-sqlite.md) | Persistência: RAM-first + SQLite | Ativo |
| ADR-001 | ~~Projeto Unity (single-repo)~~ | Substituído por ADR-004 |
| ADR-002 | ~~Workflow UI/UX Unity~~ | Substituído por ADR-004 |

---

> **Nota:** Os ADRs 007-012 (decisões de narrativa e mecânicas) estão no repositório [eras-do-brasil-gdd](https://github.com/sidinei-silva/eras-do-brasil-gdd). Branch: `decisoes-narrativa-mmorpg-2026-03-17`.

**Última atualização:** 2026-03-18
