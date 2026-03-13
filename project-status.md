# Status do Projeto — Eras do Brasil

**Última atualização:** 2026-03-11

## 🎯 Fase Atual

**Fase 0: Heartbeat (Servidor Go)** — Projeto passou por pivot completo (ADR-004): de RPG tático Unity Co-op P2P para **MUD Moderno** com servidor Go + cliente web. GDD inteiro atualizado para refletir arquitetura MMORPG server-authoritative. Stack: Go 1.22+ / HTML-CSS-JS / WebSocket / JSON. Pronto para iniciar código do servidor.

> 🔴 **Retomando?** Comece pelo [Guia de Retomada](../docs/GUIA_RETOMADA.md)
> 📋 **Decisão do Pivot:** [ADR-004](decisions/ADR-004-pivot-mmorpg-servidor-go.md)

## 📊 Progresso Geral

| Área | Status | Progresso |
|------|--------|-----------|
| **GDD (Design)** | ✅ Completo + Atualizado | 100% (atualizado para servidor Go pós-pivot) |
| **Auditoria (PR #1)** | ✅ Mergeado | 10/12 correções aplicadas, 2 recomendações futuras |
| **Decisões Técnicas** | ✅ Atualizadas | ADR-004 (pivot Go) · ADR-003 (repos, atualizado) · ADR-001/002 (substituídos) |
| **Pivot MMORPG** | ✅ Completo | GDD, Roadmap, README, docs/product, docs/tech — tudo alinhado |
| **Livros Auxiliares** | ⚠️ Parcial | 2/7 (Atlas do Eco + Catálogo de Itens) |
| **Schemas e Dados** | ✅ Completo | 100% |
| **Roadmap** | ✅ Atualizado | 7 fases (Heartbeat → Multiplayer) |
| **Plano de Desenvolvimento** | ✅ Atualizado | Calendário Mar-Dez 2026, ~6h/semana |
| **MVP Spec** | ✅ Atualizado | Corresponde à Fase 3 (Player) |
| **Comandos Admin (Dev/Teste)** | 🟡 Planejado | Trilha transversal definida em ROADMAP/backlog/Guia de Retomada |
| **Servidor Go** | 🔲 Não iniciado | Fase 0 pendente |
| **Cliente Web** | 🔲 Não iniciado | Fase 2 pendente |
| **MVP** | 🔲 Aguardando Fase 3 | 0% |

## 📕 Design Concluído

- ✅ Livro de Regras (9 capítulos — atualizado para mundo persistente/servidor)
- ✅ Livro de Classes (12 classes Tier 1 + sistema de Tiers)
- ✅ Enredo e Mundo (Ato 1 + 18 mini-campanhas + sistema Eco/Raiz)
- ✅ Design Visual (estilo, paleta, tipografia, UI Fase 1 e 2)
- ✅ Atlas do Eco Ato 1 (mapa de nós com distâncias)
- ✅ Schemas JSON (ClasseDePersonagem, Item, Inimigo)
- ✅ Game Pitch e Project Plan (atualizados para MUD Moderno)
- ✅ Auditoria do GDD (PR #1 mergeado)

## 🔥 Foco Atual — Fases de Produção

| Fase | Foco | Status |
|------|------|--------|
| **Fase 0** | Heartbeat — tick global, WebSocket, `main.go` | 🟡 Próxima |
| **Fase 1** | Mundo Vivo — NPCs, rotinas, Utility AI, fofoca | 🔲 Pendente |
| **Fase 2** | Observador — cliente web read-only, mapa, log | 🔲 Pendente |
| **Fase 3** | Jogador — criação, navegação, combate, inventário (≈ MVP) | 🔲 Pendente |
| **Fase 4** | Interação — diálogos, quests, crafting, comércio | 🔲 Pendente |
| **Fase 5** | D20 Completo — Tiers, herança, 12 classes, grid tático | 🔲 Pendente |
| **Fase 6** | Multiplayer — full loot, expedições, eventos globais | 🔲 Pendente |

> Detalhes completos em [ROADMAP.md](../ROADMAP.md) e [Guia de Retomada](../docs/GUIA_RETOMADA.md)

## 📋 Próximos Marcos

| Marco | Descrição | Dependências |
|-------|-----------|--------------|
| Fase 0 completa | Servidor Go com tick + WebSocket | Nenhuma |
| Fase 1 completa | Mundo vivo simulando NPCs 24/7 | Fase 0 |
| Fase 2 completa | Cliente web observando mundo em tempo real | Fase 1 |
| **Fase 3 = MVP** | Loop jogável: criar → explorar → lutar → coletar | Fase 2 |

## ⚠️ Bloqueios Conhecidos

- Nenhum bloqueio ativo — pronto para iniciar Fase 0 (código Go)

## 🎉 Conquistas Recentes

- ✅ **Pivot MMORPG completo** — Todo o GDD, docs, roadmap e README atualizados (ADR-004, sessão 2026-03-11)
- ✅ **Tag de segurança** `pre-pivot-mmorpg` criada antes das mudanças
- ✅ **Cap 8 — Mundo Vivo** reescrito com seções de Mundo Persistente, Tick Comutável, StoryManager, Full Loot, Arquitetura do Motor
- ✅ **Sistema Eco/Raiz** (offline pessoal + online persistente) formalizado nos Conceitos Centrais
- ✅ **Full Loot** formalizado pela primeira vez (nunca tinha sido documentado)
- ✅ **ADR-001/002** marcados como substituídos; **ADR-003** atualizado com estrutura `server/`
- ✅ **Trilha de Comandos Admin** planejada no roadmap (arquitetura, governança e entregas por fase)
- ✅ **docs/tech/** Unity docs arquivados com aviso
- ✅ **docs/product/** plano de dev, MVP spec, perfil dev — todos atualizados para Go
- ✅ **Plano de spin-offs** arquivado (decisão do pivot)
- ✅ GDD completo (9+ livros, 50+ documentos)
- ✅ Auditoria do GDD concluída (PR #1 mergeado)
- ✅ Correções pré-backlog concluídas (slots, Defesa Base, XP, Origem, pré-requisitos)
- ✅ Website atualizado com seção de notícias (PR #82 mergeado)
- ✅ Pesquisa de mercado concluída

## 📝 Notas

- **Stack:** Go 1.22+ (servidor) + HTML/CSS/JS (cliente) + WebSocket (gorilla/websocket)
- **Arquitetura:** Server-authoritative MMORPG, tick global proativo
- **Filosofia:** desenvolvimento incremental por fases, peça por peça
- **IA:** centralizada no GitHub Copilot (VS Code)
- **Decisões ativas:** [ADR-003](decisions/ADR-003-estrategia-repositorios.md) (repos) · [ADR-004](decisions/ADR-004-pivot-mmorpg-servidor-go.md) (pivot Go)
- **Decisões substituídas:** ADR-001 (Unity único) · ADR-002 (Workflow UI/UX Unity)
- **Tier 2 das classes NÃO deve ser feito agora** — é do Ato 2, após MVP

---
Este arquivo é mantido pelo @vibe-keeper. Consulte para ter um snapshot do estado atual do projeto.
