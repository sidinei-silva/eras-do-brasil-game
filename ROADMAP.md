# 🗺️ Roadmap — Eras do Brasil

> Roadmap público do projeto. Para detalhes técnicos com entregas por fase, consulte o [Roadmap de Desenvolvimento](gdd/99_Meta_e_Backlog/02_Roadmap_Desenvolvimento.md).
>
> 📐 **Decisão de pivot:** [ADR-004 — Pivot MMORPG Servidor Go](vibe/decisions/ADR-004-pivot-mmorpg-servidor-go.md)

---

## Status Atual: Fase 0 — Heartbeat (Servidor Go)

---

## Visão Geral das Fases

| Fase | Nome | Tipo | Objetivo |
|------|------|------|----------|
| **0** | Heartbeat | 💻 Código | Servidor Go com tick loop + WebSocket + cliente mínimo |
| **1** | Living World | 💻 Código | Blocos, NPCs com Utility AI, dia/noite por ticks |
| **2** | Observer | 💻 Código | Cliente web observa o mundo vivo em tempo real |
| **3** | Player | 💻 Código | Criação de personagem, D20, combate, inventário, save |
| **4** | Interaction | 💻 Código | Diálogos, fofoca, facções, crafting, Relógio da Ruptura |
| **5** | D20 Full | 💻 Código | 12 classes, 3 Origens, grid tático, Ato 1 completo |
| **6** | Multiplayer | 💻 Código | Múltiplos jogadores, full loot, missões competitivas, eventos globais |

```
Fase 0 ──► Fase 1 ──► Fase 2 ──► Fase 3 ──► Fase 4 ──► Fase 5 ──► Fase 6
```

> Cada fase entrega valor jogável completo. Detalhes em [Roadmap de Desenvolvimento](gdd/99_Meta_e_Backlog/02_Roadmap_Desenvolvimento.md).

---

## Progresso por Área

### Design & Regras (Concluído)
- [x] Livro de Regras (9 capítulos)
- [x] Livro de Classes (12 classes + sistema de Tiers)
- [x] Enredo e Mundo (Ato 1 + mini-campanhas)
- [x] Design Visual (estilo, referências de UI)
- [x] Schemas e Dados Iniciais
- [x] Auditoria do GDD
- [x] Correções pré-backlog: Slots, Defesa Base, XP, Origem, pré-requisitos

### Desenvolvimento (Em Andamento)
- [ ] Fase 0: Heartbeat (servidor Go + WebSocket + tick)
- [ ] Fase 1: Living World (blocos + NPCs + IA)
- [ ] Fase 2: Observer (cliente web observa mundo)
- [ ] Fase 3: Player (D20 + combate + inventário)
- [ ] Fase 4: Interaction (economia + diálogos + facções)
- [ ] Fase 5: D20 Full (12 classes + grid + Ato 1)
- [ ] Fase 6: Multiplayer (full loot + missões + eventos)

---

## Trilha Transversal: Comandos Admin (Dev/Teste)

> Objetivo: acelerar debug, validação de regras e testes de mundo sem depender de UI final.

### Arquitetura Planejada (sem código)

- **AdminCommandManager (goroutine dedicada):** recebe comandos administrativos, valida permissões e despacha ações para os managers do servidor.
- **Canal de entrada por fase:**
	- Fase 0-2: console local de desenvolvimento (prioritário)
	- Fase 3-5: endpoint administrativo interno (WebSocket/HTTP interno)
	- Fase 6+: opção de acesso remoto (ex.: RCON-like), somente se necessário
- **Integração com EventBus:** comandos entram por tópico administrativo e geram eventos de resposta/resultado para observabilidade.
- **Auditoria obrigatória:** todo comando administrativo gera log com ator, tick, parâmetros, resultado e impacto.
- **Isolamento de segurança:** comandos destrutivos só em modo desenvolvimento/homologação, nunca habilitados por padrão em produção.

### Categorias de Comandos

- **Inspeção:** listar NPCs, inimigos, jogadores, blocos e estado global.
- **Diagnóstico:** localizar entidade por ID/nome, ver estado detalhado, verificar filas/eventos.
- **Controle de tempo:** avançar tick, pausar, retomar, ajustar ciclo dia/noite.
- **Mutação de mundo:** mover entidade, alterar atributos, forçar spawn/despawn, reset parcial/global.
- **Teste narrativo/econômico:** disparar eventos, ajustar reputação, injetar recursos para cenários de teste.

### Entregas por Fase (Admin)

- **Fase 0:** comandos read-only de heartbeat e estado global.
- **Fase 1:** inspeção e localização de NPCs/inimigos + ciclo dia/noite.
- **Fase 2:** integração com cliente observador para depuração visual.
- **Fase 3:** comandos de personagem/inventário/save para QA funcional.
- **Fase 4:** comandos de economia, diálogos, facções e quest state.
- **Fase 5:** comandos de combate tático e balanceamento D20.
- **Fase 6:** governança multiplayer, trilhas de auditoria e controles operacionais.

---

## Decisões Técnicas

| ADR | Decisão | Status |
|-----|---------|--------|
| [ADR-001](vibe/decisions/ADR-001-organizacao-projeto-unity.md) | ~~Projeto Unity~~ | Substituído por ADR-004 |
| [ADR-002](vibe/decisions/ADR-002-workflow-ui-ux.md) | ~~Wireframes antes de UI~~ | Substituído por ADR-004 |
| [ADR-003](vibe/decisions/ADR-003-estrategia-repositorios.md) | Estratégia monorepo | Ativo |
| [ADR-004](vibe/decisions/ADR-004-pivot-mmorpg-servidor-go.md) | Pivot MMORPG — Servidor Go + Cliente Web | **Ativo** |

---

**Última atualização:** 2026-03-07
