# Backlog — Eras do Brasil

**Última atualização:** 2026-03-11

> Este arquivo é um **índice de referência rápida**. A fonte de verdade detalhada é:
>
> 📌 **Guia de Retomada** (tarefas na ordem exata): [docs/GUIA_RETOMADA.md](../docs/GUIA_RETOMADA.md)
> 📌 **Roadmap público**: [ROADMAP.md](../ROADMAP.md)
> 📌 **Decisão do Pivot**: [ADR-004](decisions/ADR-004-pivot-mmorpg-servidor-go.md)

---

## Visão Geral das Fases (Servidor Go + Cliente Web)

| Fase | Nome | Tipo | Objetivo | Status |
|------|------|------|----------|--------|
| **0** | Heartbeat | 💻 Servidor | Servidor Go roda, tick global, 1 conexão WebSocket | 🟡 Próxima |
| **1** | Mundo Vivo | 💻 Servidor | NPCs com rotinas, Utility AI, ciclo dia/noite | 🔲 Pendente |
| **2** | Observador | 💻 Full-stack | Cliente web exibe mundo em tempo real (read-only) | 🔲 Pendente |
| **3** | Jogador | 💻 Full-stack | Criação de personagem, navegação, combate, inventário | 🔲 Pendente |
| **4** | Interação | 💻 Full-stack | Diálogos, quests, crafting, comércio, save/load | 🔲 Pendente |
| **5** | D20 Completo | 💻 Full-stack | Sistema D20 full, Tiers, herança, habilidades | 🔲 Pendente |
| **6** | Multiplayer | 💻 Full-stack | Múltiplos jogadores, full loot, expedições na Raiz | 🔲 Pendente |

---

## Documentação Concluída (Pré-Código)

- [x] GDD completo (9 capítulos do Livro de Regras)
- [x] 12 classes Tier 1 documentadas
- [x] Enredo e Mundo (Ato 1 + 18 mini-campanhas)
- [x] Design Visual (pixel art, paleta, UI Fase 1 e 2)
- [x] Schemas JSON (ClasseDePersonagem, Item, Inimigo)
- [x] Auditoria do GDD concluída (PR #1 mergeado)
- [x] Correções pré-backlog (slots, Defesa Base, XP, Origem, pré-requisitos)
- [x] **Pivot MMORPG** — GDD inteiro atualizado para servidor Go (ADR-004)

---

## Tarefas por Fase

### Workstream Transversal — Comandos Admin (Dev/Teste)

- [ ] Definir catálogo inicial de comandos administrativos (inspeção, diagnóstico, mutação, reset)
- [ ] Definir contrato de entrada/saída dos comandos (request/response/evento)
- [ ] Definir modelo de autorização por ambiente (dev/homolog/prod)
- [ ] Definir trilha de auditoria obrigatória (quem executou, quando, impacto)
- [ ] Definir estratégia de rollback para comandos destrutivos

**Arquitetura alvo (sem código):**
- [ ] 1 goroutine dedicada para AdminCommandManager
- [ ] integração com EventBus para despacho e resposta
- [ ] fonte de entrada local no início (console dev)
- [ ] endpoint administrativo interno nas fases intermediárias
- [ ] remoto opcional (RCON-like) apenas quando houver necessidade real

### Fase 0 — Heartbeat (servidor Go)
- [ ] `main.go` com tick global (`time.Ticker` + goroutine)
- [ ] Struct `Mundo` com `ProcessarTick()`
- [ ] WebSocket listener (`gorilla/websocket`)
- [ ] Cliente HTML mínimo que conecta e recebe ticks
- [ ] Admin v0: comandos read-only de status (tick, uptime, estado global)

### Fase 1 — Mundo Vivo
- [ ] NPCs com rotinas diárias (Utility AI)
- [ ] Sistema de Fofoca entre NPCs
- [ ] Ciclo dia/noite afetando comportamento
- [ ] Eventos de mundo (aleatórios + sazonais)
- [ ] StoryManager — arco narrativo dinâmico
- [ ] Admin v1: listar/localizar NPCs e inimigos, detalhes de estado por entidade

### Fase 2 — Observador
- [ ] Cliente web mostra estado do mundo em tempo real
- [ ] Mapa de nós navegável (read-only)
- [ ] Log de eventos visível
- [ ] UI básica (HTML/CSS/JS)
- [ ] Admin v2: painel/visão observável dos comandos de inspeção

### Fase 3 — Jogador (≈ MVP "O Despertar")
- [ ] Criação de personagem (1 classe: Guerreiro Tribal)
- [ ] Navegação por blocos com custo de Ticks
- [ ] Combate estático (D20 simplificado)
- [ ] Inventário e loot
- [ ] HUD principal
- [ ] Save/Load do estado do jogador
- [ ] Admin v3: comandos de personagem/inventário/save para testes funcionais

### Fase 4 — Interação
- [ ] Diálogos ramificados com NPCs
- [ ] Sistema de Quests
- [ ] Crafting e coleta com proficiências
- [ ] Comércio com NPCs (server-authoritative)
- [ ] Facções e reputação
- [ ] Status e condições em combate
- [ ] Mini-campanha "O Caçador que Não Voltou"
- [ ] Admin v4: comandos de facção, economia, diálogos e estado de quests

### Fase 5 — D20 Completo
- [ ] Sistema D20 completo (vantagem, desvantagem, críticos)
- [ ] Tiers 1→2→3 e Moedas de Classe
- [ ] Herança de habilidades (Dom da Revivência)
- [ ] Habilidades ativas em combate
- [ ] Balanceamento com 12 classes Tier 1
- [ ] Grid tático (Fase 2 visual do combate)
- [ ] Admin v5: comandos de simulação de combate e balanceamento

### Fase 6 — Multiplayer
- [ ] Múltiplas conexões simultâneas
- [ ] Full Loot (com mitigações: Marca do Eco, seguro)
- [ ] Expedições na Raiz (dungeons competitivas)
- [ ] Missões competitivas entre jogadores
- [ ] Eventos globais (rupturas, mudanças de era)
- [ ] Inimigos evolutivos (adaptam táticas)
- [ ] Economia multiplayer server-authoritative
- [ ] Admin v6: governança operacional multiplayer + auditoria avançada

---

## Após MVP — Futuro (NÃO fazer agora)

- [ ] Detalhar habilidades de Tier 2 para as 12 classes
- [ ] Criar Reaction Abilities (1 por classe)
- [ ] Criar Livro de Magias e Rituais
- [ ] Criar Livro de Receitas de Crafting

---

**Próxima ação recomendada:**
- [ ] Iniciar Fase 0 — criar `server/main.go` com tick loop e WebSocket

**Gerenciado por:** @tech-lead
**Referência:** [ROADMAP.md](../ROADMAP.md) · [Guia de Retomada](../docs/GUIA_RETOMADA.md) · [ADR-004](decisions/ADR-004-pivot-mmorpg-servidor-go.md) · [ADRs](decisions/)
