# Sessão 5 — Pivot MMORPG: Servidor Go + Cliente Web

**Data:** 2026-03-11
**Duração:** Sessão longa (múltiplas interações)
**Tipo:** Pivot de arquitetura + Rollback massivo de documentação

---

## 🎯 Objetivo

Abandonar a abordagem Unity Co-op P2P (que gerou 3 anos de paralisia por documentação) e pivotar para um **MUD Moderno** com:
- **Servidor Go** (goroutines para tick loop, NPC AI, WebSocket)
- **Cliente Web** (HTML/CSS/JS via WebSocket, zero install)
- **Arquitetura MMORPG** server-authoritative

## 💡 Motivação

O desenvolvedor expressou frustração com 3 anos de documentação sem nenhuma linha de código de jogo. A análise revelou:
- Unity/C# era uma barreira (não é a stack primária do dev)
- 38 POCs planejadas antes de qualquer jogo jogável
- Spin-offs adiavam ainda mais o projeto principal
- A visão original do projeto (repo legado `sidinei-silva/ErasDoBrasil-Historico-Legado`, 2023) era MMORPG, que foi diluída em 4 commits de novembro/dezembro 2025

## 📋 Decisões Tomadas

### ADR-004 — Pivot MMORPG Servidor Go
- **Go** sobre Node.js: goroutines para NPC AI concorrente (cada NPC = goroutine leve)
- **Go** sobre PHP: estado persistente em memória (PHP = request/response stateless)
- **Web client** sobre Unity client: zero install, skills existentes do dev (sênior web), deploy imediato
- **Server-authoritative**: anti-cheat by design, world state em memória
- **Tick global proativo**: `time.Ticker` + goroutine, `Mundo.ProcessarTick()` comutável (reativo offline / proativo online)

### Análise do Repo Legado
- 4 commits que removeram features MMO: `c4fff4a`, `820ad6e`, `16132ad`, `ed287af`
- Comparativo criado em `vibe/scratchpad/comparativo-gdd-legado-vs-atual.md`
- Conclusão: manter GDD atual (superior em Caps 4-9) e restaurar cirurgicamente conceitos MMO

### Novas Features Formalizadas
- **Full Loot** (Cap. 8, §8.10) — nunca tinha sido documentado formalmente
- **Sistema Eco/Raiz** (Cap. 00, §4) — Eco = offline pessoal, Raiz = online persistente
- **StoryManager** (Cap. 8, §8.7) — arco narrativo dinâmico server-side
- **Missões Competitivas** (Cap. 8, §8.8) — jogadores competem pelo mesmo objetivo
- **Eventos Globais** (Cap. 8, §8.9) — rupturas e mudanças de era

## 📝 Arquivos Modificados

### GDD — Livro de Regras
| Arquivo | Mudança |
|---------|---------|
| `gdd/01_Livro_de_Regras/01_Introducao_e_Ambientacao.md` | RPG Digital → MUD Moderno |
| `gdd/01_Livro_de_Regras/04_Sistema_de_Combate.md` | Persistente |
| `gdd/01_Livro_de_Regras/05_Exploracao_e_Mundo.md` | 6 replacements (persistente, server, events) |
| `gdd/01_Livro_de_Regras/06_Economia_Itens_e_Crafting.md` | 4 replacements (persistente, ticks, server trade) |
| `gdd/01_Livro_de_Regras/07_Magia_e_Espiritualidade.md` | 1 replacement |
| `gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md` | **MAJOR REWRITE** — §8.1-8.2 + §8.7-8.12 novos |

### GDD — Enredo e Mundo
| Arquivo | Mudança |
|---------|---------|
| `gdd/03_Enredo_e_Mundo/00_Conceitos_Centrais_do_Mundo.md` | Co-op → Eco/Raiz (sistema híbrido) |
| `gdd/03_Enredo_e_Mundo/01_Ato_1_A_Primeira_Ruptura.md` | Co-op → Modo Online (Raiz) |

### GDD — Meta e Backlog
| Arquivo | Mudança |
|---------|---------|
| `gdd/Project Plan.md` | **REWRITE COMPLETO** — MUD Moderno |
| `gdd/Game_Pitch.md` | Pillar Co-op → Mundo Persistente, §3.4 reescrito |
| `gdd/99_Meta_e_Backlog/01_Estrategia_de_Dados_e_Documentacao.md` | Pipeline → server/data/ |
| `gdd/99_Meta_e_Backlog/02_Roadmap_Desenvolvimento.md` | **REWRITE MASSIVO** — Stack/Fases/POCs removidas |
| `gdd/README.md` | Pilares, links, descrição atualizados |

### Raiz do Projeto
| Arquivo | Mudança |
|---------|---------|
| `README.md` | **REWRITE COMPLETO** — MUD Moderno, stack, fases |
| `ROADMAP.md` | **REWRITE COMPLETO** — 7 fases, ADR-004 |

### docs/product/
| Arquivo | Mudança |
|---------|---------|
| `docs/product/mvp-o-despertar-spec.md` | Stack Go, referências de fase, critério WebSocket |
| `docs/product/perfil-developer.md` | Skills Go, removeu Unity |
| `docs/product/plano-spinoffs-aprendizado.md` | ⚠️ ARQUIVADO |
| `docs/product/plano-desenvolvimento-realista.md` | **REWRITE COMPLETO** — Calendário Go |

### docs/tech/
| Arquivo | Mudança |
|---------|---------|
| `docs/tech/organizacao-projeto-unity.md` | ⚠️ ARQUIVADO |
| `docs/tech/workflow-ui-ux.md` | ⚠️ ARQUIVADO |

### docs/
| Arquivo | Mudança |
|---------|---------|
| `docs/GUIA_RETOMADA.md` | **REWRITE COMPLETO** — ~740 linhas reescritas para Fases Go |

### vibe/decisions/
| Arquivo | Mudança |
|---------|---------|
| `vibe/decisions/ADR-001-organizacao-projeto-unity.md` | Status: Substituído por ADR-004 |
| `vibe/decisions/ADR-002-workflow-ui-ux.md` | Status: Substituído por ADR-004 |
| `vibe/decisions/ADR-003-estrategia-repositorios.md` | Aceito + atualizado (server/, removidos game/pocs/spinoffs) |
| `vibe/decisions/ADR-004-pivot-mmorpg-servidor-go.md` | **CRIADO** — Decisão completa do pivot |

### vibe/
| Arquivo | Mudança |
|---------|---------|
| `vibe/backlog.md` | **REWRITE** — Sprints → Fases, tarefas por fase |
| `vibe/project-status.md` | **REWRITE** — Status atualizado pós-pivot |

### pocs/
| Arquivo | Mudança |
|---------|---------|
| `pocs/README.md` | ⚠️ ARQUIVO HISTÓRICO — aviso de ADR-004 |

## 📊 Resumo Numérico

- **Arquivos modificados:** ~30
- **Rewrites completos:** 8 (Project Plan, Game Pitch §3.4, Roadmap Dev, ROADMAP, README, Guia Retomada, backlog, project-status)
- **Arquivos arquivados:** 4 (Unity org, UI/UX workflow, spin-offs, POCs README)
- **ADRs:** 1 criado (004), 2 substituídos (001/002), 1 atualizado (003)
- **Novas seções do GDD:** 6 (Full Loot, Eco/Raiz, StoryManager, Missões Competitivas, Eventos Globais, Arquitetura do Motor)
- **Tag de segurança:** `pre-pivot-mmorpg`

## 🔮 Próximo Passo

Iniciar **Fase 0 — Heartbeat**: criar `server/main.go` com tick global + WebSocket.

---

**Gerenciado por:** @vibe-keeper
