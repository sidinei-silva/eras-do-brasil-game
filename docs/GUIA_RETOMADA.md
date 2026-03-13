# 🧭 Guia de Retomada — Eras do Brasil

> **Objetivo:** Este documento é o seu ponto de partida para retomar o projeto. Ele organiza **tudo** que precisa ser feito, na **ordem exata** em que deve ser feito, com links para cada documento relevante.
>
> **Como usar:** Siga as Fases na ordem. Cada Fase tem entregas numeradas. Não pule Fases.
>
> **Stack:** Go 1.22+ (servidor) · HTML/CSS/JS (cliente) · WebSocket (gorilla/websocket) · JSON
>
> **Decisão do Pivot:** [ADR-004](../vibe/decisions/ADR-004-pivot-mmorpg-servidor-go.md)
>
> **Última atualização:** 2026-03-11

---

## 📍 Onde Você Está Agora

Antes de qualquer tarefa, entenda o estado atual:

| Área | Status | O que existe |
|------|--------|-------------|
| **GDD (Design)** | ✅ 100% Completo | 9 capítulos de regras, 12 classes Tier 1, 18 mini-campanhas, sistema Eco/Raiz |
| **Auditoria (PR #1)** | ✅ Mergeado | 10/12 correções aplicadas ([ver auditoria](../gdd/99_Meta_e_Backlog/03_Auditoria_GDD.md)) |
| **Pivot MMORPG** | ✅ Completo | Todo o GDD atualizado para servidor Go + cliente web ([ADR-004](../vibe/decisions/ADR-004-pivot-mmorpg-servidor-go.md)) |
| **Decisões Técnicas** | ✅ Atualizadas | ADR-004 (pivot Go) · ADR-003 (repos) · ADR-001/002 (substituídos) |
| **Livros Auxiliares** | ⚠️ 2/7 | Atlas do Eco + Livro de Itens existem. 5 livros referenciados no GDD ainda não existem |
| **Código** | 🔲 Não iniciado | Servidor Go não criado, cliente web não criado |

### O Que Foi Feito Antes do Código

Estas são conquistas de documentação **já concluídas**, não precisam ser refeitas:

- ✅ GDD completo com 9 capítulos de regras + 12 classes Tier 1
- ✅ Auditoria do GDD (PR #1 mergeado — 4 correções, 7 recomendações)
- ✅ Correções pré-backlog: Slots de habilidades, Defesa Base, XP/Níveis, Troca de Origem, Pré-requisitos
- ✅ Schemas JSON (ClasseDePersonagem, Item, Inimigo) em `gdd/06_Dados_e_Assets/`
- ✅ Atlas do Eco Ato 1 (mapa de nós com distâncias)
- ✅ Design Visual completo (estilo, paleta, UI Fase 1 e 2)
- ✅ Pivot completo de Unity Co-op P2P → Servidor Go MMORPG (ADR-004)
- ✅ Full Loot formalizado no GDD (Cap. 8, seção 8.10)
- ✅ Sistema Eco/Raiz (offline + online) definido nos Conceitos Centrais

---

## 🗓️ Visão Geral das Fases

> Cada Fase é incremental — o servidor vai ganhando capacidades até se tornar um jogo completo.
> Roadmap completo em [ROADMAP.md](../ROADMAP.md)

| Fase | Nome | Tipo | O que entrega |
|------|------|------|---------------|
| **0** | [Heartbeat](#-fase-0--heartbeat) | 💻 Servidor | Tick global + 1 WebSocket |
| **1** | [Mundo Vivo](#-fase-1--mundo-vivo) | 💻 Servidor | NPCs com rotinas e IA |
| **2** | [Observador](#-fase-2--observador) | 💻 Full-stack | Cliente web read-only |
| **3** | [Jogador (≈ MVP)](#-fase-3--jogador--mvp-o-despertar) | 💻 Full-stack | Criar, explorar, lutar, coletar |
| **4** | [Interação](#-fase-4--interação) | 💻 Full-stack | Diálogos, quests, crafting, comércio |
| **5** | [D20 Completo](#-fase-5--d20-completo) | 💻 Full-stack | Sistema completo, 12 classes, Tiers |
| **6+** | [Multiplayer](#-fase-6--multiplayer) | 💻 Full-stack | Full loot, expedições, eventos globais |

```
Fase 0 (Heartbeat) ──► Fase 1 (Mundo Vivo) ──► Fase 2 (Observador)
                                                      │
                              Fase 3 (Jogador/MVP) ◄──┘
                                      │
                              Fase 4 (Interação)
                                      │
                              Fase 5 (D20 Completo)
                                      │
                              Fase 6 (Multiplayer)
```

---

## 🛠️ Trilha Transversal — Comandos Admin (Dev/Teste)

> Esta trilha acompanha todas as fases para acelerar debug, inspeção de estado e testes de regressão sem depender da UI final.

### Arquitetura Planejada (sem código)

- **AdminCommandManager (goroutine dedicada):** recebe comandos administrativos e coordena execução.
- **Fluxo de comunicação:** Entrada admin -> validação/permissão -> despacho para manager de domínio -> resposta -> auditoria.
- **Integração com EventBus:** comandos e respostas circulam como eventos administrativos.
- **Canais por maturidade:**
        - Fase 0-2: console local (prioritário)
        - Fase 3-5: endpoint interno administrativo (HTTP/WS interno)
        - Fase 6+: remoto opcional (RCON-like), somente se houver necessidade real
- **Governança:** comando destrutivo só em dev/homolog e sempre com log de auditoria.

### Entregas Admin por Fase

- **Fase 0:** status global, tick atual, uptime e saúde do loop.
- **Fase 1:** listar/localizar NPCs e inimigos, inspecionar estado por entidade.
- **Fase 2:** observabilidade dos comandos no cliente observador.
- **Fase 3:** inspeção de jogador, inventário e fluxo de save/load.
- **Fase 4:** comandos de facção, economia, diálogos e estado de quests.
- **Fase 5:** comandos de simulação para balanceamento D20 e combate.
- **Fase 6:** governança operacional multiplayer com trilha avançada de auditoria.

---

## 💻 Fase 0 — Heartbeat

> **Objetivo:** Provar que o servidor Go funciona — tick global rodando, uma conexão WebSocket recebendo estado.
>
> **Estimativa:** 1–2 sessões
> **Pré-requisito:** Go 1.22+ instalado

### Estrutura de Pastas do Servidor

```
server/
├── main.go             (entry point + tick loop)
├── world/              (Mundo, Bloco, NPC)
├── combat/             (Motor D20, turnos, status)
├── economy/            (Inventário, crafting, comércio)
├── data/               (JSONs, structs de dados)
└── go.mod
```

**Referência:** [ADR-003 — Estrutura de Repositórios](../vibe/decisions/ADR-003-estrategia-repositorios.md)

### Entregas

| # | Entrega | Critério de Aceite |
|---|---------|-------------------|
| 0.1 | `go mod init` + estrutura de pastas | Compilar sem erros |
| 0.2 | `main.go` com tick global | `time.Ticker` + goroutine imprimindo "Tick N" no console |
| 0.3 | Struct `Mundo` com `ProcessarTick()` | Struct criada em `world/`, chamada pelo tick loop |
| 0.4 | WebSocket listener | Cliente HTML conecta, recebe JSON `{"tick": N}` a cada tick |

### Referências do GDD

- **Tick Comutável:** [`gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md`](../gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md) §8.2
- **Arquitetura do Motor:** [`gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md`](../gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md) §8.12

### ✅ Checklist da Fase 0

- [ ] 0.1 — `go mod init` + pastas (`server/`, `world/`, `combat/`, `economy/`, `data/`)
- [ ] 0.2 — `main.go` com `time.Ticker` + goroutine de tick global
- [ ] 0.3 — Struct `Mundo` em `world/mundo.go` com `ProcessarTick()`
- [ ] 0.4 — WebSocket listener (`gorilla/websocket`) enviando tick ao cliente

---

## 💻 Fase 1 — Mundo Vivo

> **Objetivo:** O servidor simula um mundo vivo — NPCs existem, têm rotinas, tomam decisões. O mundo funciona sem jogadores.
>
> **Estimativa:** 2–4 sessões
> **Pré-requisito:** Fase 0

### Entregas

| # | Entrega | Critério de Aceite |
|---|---------|-------------------|
| 1.1 | Structs de NPC com rotinas | NPC segue agenda diária (acordar, trabalhar, comer, dormir) |
| 1.2 | Utility AI básica | NPC escolhe ação com base em necessidades (fome, fadiga, medo) |
| 1.3 | Ciclo dia/noite | Tick avança relógio, NPCs mudam comportamento à noite |
| 1.4 | Sistema de Fofoca | NPCs compartilham informações sobre eventos |
| 1.5 | StoryManager (sementes) | Arco narrativo reage a ações no mundo |

### Referências do GDD

- **Utility AI:** [`gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md`](../gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md) §8.4
- **Rotinas de NPCs:** [`gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md`](../gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md) §8.3
- **Fofoca:** [`gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md`](../gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md) §8.5
- **StoryManager:** [`gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md`](../gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md) §8.7

### ✅ Checklist da Fase 1

- [ ] 1.1 — Structs NPC em `world/npc.go` com rotina diária
- [ ] 1.2 — Utility AI: Score = Peso × (1 – Necessidade Normalizada)
- [ ] 1.3 — Ciclo dia/noite vinculado ao tick global
- [ ] 1.4 — NPCs propagam informações entre si via fofoca
- [ ] 1.5 — StoryManager: sementes narrativas reagem a threshold de eventos

---

## 💻 Fase 2 — Observador

> **Objetivo:** Um cliente web se conecta ao servidor e mostra o mundo em tempo real. O jogador ainda não interage — apenas observa.
>
> **Estimativa:** 2–3 sessões
> **Pré-requisito:** Fase 1

### Entregas

| # | Entrega | Critério de Aceite |
|---|---------|-------------------|
| 2.1 | Cliente HTML/CSS/JS básico | Página conecta via WebSocket, exibe dados |
| 2.2 | Mapa de nós (read-only) | Mapa mostra blocos com NPCs se movendo |
| 2.3 | Log de eventos | Painel lateral com eventos do mundo em tempo real |
| 2.4 | HUD de tempo | Relógio de ticks, ciclo dia/noite, clima |

### Referências do GDD

- **Mapa de Nós:** [`gdd/04_Design_Visual/05_UI_Fase_1_Exploracao_e_Combate.md`](../gdd/04_Design_Visual/05_UI_Fase_1_Exploracao_e_Combate.md)
- **Atlas do Eco:** [`gdd/05_Livros_Auxiliares/01_Atlas_do_Eco_Ato1.md`](../gdd/05_Livros_Auxiliares/01_Atlas_do_Eco_Ato1.md)

### ✅ Checklist da Fase 2

- [ ] 2.1 — `web/index.html` conecta ao servidor e exibe JSON recebido
- [ ] 2.2 — Mapa de nós renderizado mostrando posição de NPCs
- [ ] 2.3 — Log de eventos em painel lateral (lista rolável)
- [ ] 2.4 — Relógio de ticks + indicador dia/noite

---

## 💻 Fase 3 — Jogador (≈ MVP "O Despertar")

> **Objetivo:** O jogador pode criar um personagem e jogar. Loop completo de 15–30 minutos: criar → explorar → lutar → coletar → descansar.
>
> **Estimativa:** 4–6 sessões
> **Pré-requisito:** Fase 2
>
> 📌 **Esta Fase corresponde ao MVP.** Veja a spec completa em [`docs/product/mvp-o-despertar-spec.md`](product/mvp-o-despertar-spec.md)

### O Core Loop do MVP

```
Criar Personagem (Origem Indígena → Guerreiro Tribal)
    → Explorar Bloco (Vila de São Tomé → Floresta do Norte → Ruínas Queimadas)
    → Encontrar Combate (Lobo / Espírito Menor)
    → Rolar D20, usar habilidades, receber/causar dano
    → Coletar Loot/XP
    → Voltar à Vila (descansar, equipar)
    → Repetir
```

### Entregas

| # | Entrega | Critério de Aceite |
|---|---------|-------------------|
| 3.1 | Criação de personagem | Escolher origem + classe, point-buy (27 pts), stats calculados |
| 3.2 | Navegação por blocos | Gastar Ticks para se mover, custo por terreno |
| 3.3 | Combate estático (D20 simplificado) | Iniciativa → Turnos → D20 vs Defesa → Dano → Loot |
| 3.4 | Inventário e equipamentos | Equipar/desequipar, peso/capacidade |
| 3.5 | HUD principal | PV, XP, recursos, relógio, posição |
| 3.6 | Save/Load | Estado do jogador persiste entre sessões |

### Escopo Exato do MVP

| Elemento | Quantidade |
|----------|-----------|
| Classes jogáveis | 1 (Guerreiro Tribal) |
| Blocos do mapa | 3 (Vila de São Tomé, Floresta do Norte, Ruínas Queimadas) |
| Inimigos | 2 (Lobo Cinza, Espírito Menor) |
| NPCs interativos | 2–3 (Vila de São Tomé) |
| Itens | 10–15 (armas, armaduras, consumíveis, recursos) |
| Quest | 1 ("O Caçador que Não Voltou" — simplificada) |

### Referências do GDD

- **Criação:** [`gdd/01_Livro_de_Regras/03_Criacao_de_Personagem.md`](../gdd/01_Livro_de_Regras/03_Criacao_de_Personagem.md)
- **Combate:** [`gdd/01_Livro_de_Regras/04_Sistema_de_Combate.md`](../gdd/01_Livro_de_Regras/04_Sistema_de_Combate.md)
- **Navegação:** [`gdd/05_Livros_Auxiliares/01_Atlas_do_Eco_Ato1.md`](../gdd/05_Livros_Auxiliares/01_Atlas_do_Eco_Ato1.md)
- **Economia/Itens:** [`gdd/01_Livro_de_Regras/06_Economia_Itens_e_Crafting.md`](../gdd/01_Livro_de_Regras/06_Economia_Itens_e_Crafting.md)
- **Guerreiro Tribal:** [`gdd/02_Livro_de_Classes/2_Origem_Indigena/Guerreiro_Tribal.md`](../gdd/02_Livro_de_Classes/2_Origem_Indigena/Guerreiro_Tribal.md)
- **Dados mockup:** [`gdd/06_Dados_e_Assets/03_Dados_Iniciais_Ato1.md`](../gdd/06_Dados_e_Assets/03_Dados_Iniciais_Ato1.md)

### ✅ Checklist da Fase 3

- [ ] 3.1 — Criação de personagem (point-buy, 1 classe, stats)
- [ ] 3.2 — Navegação por blocos (Vila→Floresta→Ruínas, custo de Ticks)
- [ ] 3.3 — Combate estático (D20 simplificado, Iniciativa, Turnos, Loot)
- [ ] 3.4 — Inventário (equipar, peso, capacidade)
- [ ] 3.5 — HUD principal (PV, XP, relógio, posição)
- [ ] 3.6 — Save/Load (JSON)
- [ ] Playtest: jogar loop completo 3+ vezes, ajustar números

---

## 💻 Fase 4 — Interação

> **Objetivo:** O mundo reage ao jogador. Diálogos, quests, crafting e comércio.
>
> **Estimativa:** 4–6 sessões
> **Pré-requisito:** Fase 3

### Entregas

| # | Entrega | Critério de Aceite |
|---|---------|-------------------|
| 4.1 | Diálogos ramificados com NPCs | Árvore de diálogo com condições (reputação, itens) |
| 4.2 | Sistema de Quests | Aceitar, rastrear progresso, completar, recompensa |
| 4.3 | Crafting e coleta | Proficiências de vida, coleta de recursos, receitas |
| 4.4 | Comércio com NPCs | Comprar/vender (server-authoritative) |
| 4.5 | Facções e reputação | Ações afetam reputação, NPCs reagem |
| 4.6 | Status e condições em combate | Envenenado, Atordoado, Queimando — efeitos por turno |
| 4.7 | Mini-campanha completa | "O Caçador que Não Voltou" jogável do início ao fim |

### Referências do GDD

- **Exploração/Mundo:** [`gdd/01_Livro_de_Regras/05_Exploracao_e_Mundo.md`](../gdd/01_Livro_de_Regras/05_Exploracao_e_Mundo.md)
- **Economia/Crafting:** [`gdd/01_Livro_de_Regras/06_Economia_Itens_e_Crafting.md`](../gdd/01_Livro_de_Regras/06_Economia_Itens_e_Crafting.md)
- **Facções:** [`gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md`](../gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md) §8.6
- **Mini-campanha:** [`gdd/03_Enredo_e_Mundo/1_MiniCampanhas_Indigena/`](../gdd/03_Enredo_e_Mundo/1_MiniCampanhas_Indigena/)

### ✅ Checklist da Fase 4

- [ ] 4.1 — Diálogos ramificados (árvore de diálogo JSON)
- [ ] 4.2 — Sistema de Quests (aceitar, tracker, completar)
- [ ] 4.3 — Crafting e coleta (proficiências, recursos, receitas)
- [ ] 4.4 — Comércio com NPCs (server-authoritative)
- [ ] 4.5 — Facções e reputação
- [ ] 4.6 — Status/condições em combate
- [ ] 4.7 — Mini-campanha "O Caçador que Não Voltou" jogável

---

## 💻 Fase 5 — D20 Completo

> **Objetivo:** Sistema D20 completo com todas as classes, Tiers e grid tático.
>
> **Estimativa:** 4–6 sessões
> **Pré-requisito:** Fase 4

### Entregas

| # | Entrega | Critério de Aceite |
|---|---------|-------------------|
| 5.1 | Sistema D20 completo | Vantagem, desvantagem, críticos, todas as regras |
| 5.2 | Tiers 1→2→3 | Moedas de Classe, evolução, pré-requisitos |
| 5.3 | Herança de habilidades | Dom da Revivência: trocar classe mantendo habilidades |
| 5.4 | Habilidades ativas em combate | Custo de recurso, efeitos de área, recarga |
| 5.5 | 12 classes Tier 1 jogáveis | Balanceamento completo |
| 5.6 | Grid tático (Fase 2 visual) | Combate com posicionamento e AoE |

### Referências do GDD

- **D20:** [`gdd/01_Livro_de_Regras/02_Mecanicas_Basicas.md`](../gdd/01_Livro_de_Regras/02_Mecanicas_Basicas.md)
- **Tiers:** [`gdd/02_Livro_de_Classes/01_Tiers_e_Evolucao.md`](../gdd/02_Livro_de_Classes/01_Tiers_e_Evolucao.md)
- **Herança:** [`gdd/02_Livro_de_Classes/03_Heranca_de_Habilidades.md`](../gdd/02_Livro_de_Classes/03_Heranca_de_Habilidades.md)
- **UI Combate Tático:** [`gdd/04_Design_Visual/06_UI_Fase_2_Combate_Tatico.md`](../gdd/04_Design_Visual/06_UI_Fase_2_Combate_Tatico.md)

### ✅ Checklist da Fase 5

- [ ] 5.1 — D20 completo (vantagem, críticos, regras)
- [ ] 5.2 — Tiers 1→2→3 (Moedas de Classe, evolução)
- [ ] 5.3 — Herança de habilidades (Dom da Revivência)
- [ ] 5.4 — Habilidades ativas em combate
- [ ] 5.5 — 12 classes Tier 1 balanceadas
- [ ] 5.6 — Grid tático (posicionamento, AoE)

---

## 💻 Fase 6+ — Multiplayer

> **Objetivo:** Múltiplos jogadores no mesmo mundo persistente. Full loot, expedições competitivas, eventos globais.
>
> **Estimativa:** Futuro (após Fase 5 estável)
> **Pré-requisito:** Fase 5
>
> ⚠️ **Esta Fase só deve ser iniciada após o jogo single-player estar jogável e balanceado.**

### Entregas

| # | Entrega | Critério de Aceite |
|---|---------|-------------------|
| 6.1 | Múltiplas conexões | 2+ jogadores no mesmo mundo simultaneamente |
| 6.2 | Full Loot | Morte = perda de inventário (com Marca do Eco e seguro) |
| 6.3 | Expedições na Raiz | Dungeons competitivas para grupos |
| 6.4 | Missões competitivas | Jogadores competem pelo mesmo objetivo |
| 6.5 | Eventos globais | Rupturas, mudanças de era, afetam todos |
| 6.6 | Inimigos evolutivos | NPCs adaptam táticas com base em padrões dos jogadores |
| 6.7 | Economia multiplayer | Comércio entre jogadores (server-authoritative) |

### Referências do GDD

- **Full Loot:** [`gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md`](../gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md) §8.10
- **Missões Competitivas:** [`gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md`](../gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md) §8.8
- **Eventos Globais:** [`gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md`](../gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md) §8.9
- **Eco/Raiz:** [`gdd/03_Enredo_e_Mundo/00_Conceitos_Centrais_do_Mundo.md`](../gdd/03_Enredo_e_Mundo/00_Conceitos_Centrais_do_Mundo.md) §4

### ✅ Checklist da Fase 6

- [ ] 6.1 — Múltiplas conexões WebSocket simultâneas
- [ ] 6.2 — Full Loot (Marca do Eco, seguro, mitigações)
- [ ] 6.3 — Expedições na Raiz (dungeons competitivas)
- [ ] 6.4 — Missões competitivas entre jogadores
- [ ] 6.5 — Eventos globais (rupturas, era shifts)
- [ ] 6.6 — Inimigos evolutivos
- [ ] 6.7 — Economia multiplayer server-authoritative

---

## 📝 Documentação Pendente (Livros Auxiliares)

> Estes livros são referenciados no GDD mas ainda não foram criados. Crie-os conforme necessário em cada Fase.

| Livro | Necessário na Fase | Quando Criar |
|-------|:------------------:|-------------|
| **Livro de Habilidades** | 3 | Consolidar habilidades das 12 classes num índice |
| **Livro de Proficiências de Vida** | 4 | Coleta, produção, níveis, XP por proficiência |
| **Bestiário Ato 1** | 3 | Stats, lore, loot, comportamento IA dos inimigos |
| **Livro de Magias e Rituais** | 5 | Sistema mágico completo (Origem Folclórica e Indígena) |
| **Livro de Receitas de Crafting** | 4 | Receitas + ingredientes + qualidade do resultado |

### Livros que JÁ existem:
- ✅ [Atlas do Eco Ato 1](../gdd/05_Livros_Auxiliares/01_Atlas_do_Eco_Ato1.md) — Mapa de nós com distâncias
- ✅ [Livro de Itens e Equipamentos](../gdd/05_Livros_Auxiliares/03_Livro_de_Itens_e_Equipamentos.md) — 15 itens MVP, Matriz 5×5

---

## 📝 Lore Pendente (da Auditoria)

> Issues de lore identificadas na [Auditoria do GDD](../gdd/99_Meta_e_Backlog/03_Auditoria_GDD.md). Resolva conforme surgirem durante o desenvolvimento.

| # | Questão | Impacto |
|---|---------|---------|
| 1 | Dom da Revivência: latente ou ativo no Ato 1? | Afeta narrativa da Fase 3+ |
| 2 | Causa da Ruptura de 1497: definir origem canônica | Lore central |
| 3 | Motivação da facção Folclórica no Ato 1 | Afeta quests Fase 4 |
| 4 | "Mana" referenciada sem sistema — Fadiga vs Mana | Afeta classes mágicas Fase 5 |
| 5 | Dano Espiritual / Exaustão Espiritual não definidos | Combate Fase 5 |
| 6 | Escala de Ticks: 500 Ticks = quanto tempo real? | Fase 0 (definir ao implementar) |
| 7 | Guardião da Fenda: identidade e stats | Fase 4+ (boss) |
| 8 | Sanidade/Moral: definir ou remover | Fase 4+ |
| 9 | "Ecoera": definir ou remover | Lore |
| 10 | 12 mini-campanhas sem impacto sinérgico no Ato 1 | Fase 4+ |

---

## 📚 Mapa de Navegação do Projeto

> Para não se perder, use este mapa mental de onde encontrar cada coisa:

```
eras-do-brasil/
├── 📄 README.md                    ← Visão geral do projeto
├── 📄 ROADMAP.md                   ← Status geral e fases
├── 📄 docs/GUIA_RETOMADA.md        ← VOCÊ ESTÁ AQUI 🔴
│
├── 📂 gdd/                         ← Game Design Document (a Bíblia)
│   ├── 📄 README.md                ← Índice do GDD com links
│   ├── 📄 Game_Pitch.md            ← Apresentação do jogo
│   ├── 📄 Project Plan.md          ← Plano de projeto (MUD Moderno)
│   │
│   ├── 📂 01_Livro_de_Regras/      ← 9 capítulos de mecânicas
│   ├── 📂 02_Livro_de_Classes/     ← 12 classes Tier 1 + sistema
│   ├── 📂 03_Enredo_e_Mundo/       ← Ato 1 + 18 mini-campanhas + Eco/Raiz
│   ├── 📂 04_Design_Visual/        ← Pixel art, UI, HUD (specs de referência)
│   ├── 📂 05_Livros_Auxiliares/    ← Atlas do Eco + futuros livros
│   ├── 📂 06_Dados_e_Assets/       ← Schemas JSON, dados mockup
│   └── 📂 99_Meta_e_Backlog/       ← Estratégia, Roadmap, Auditoria
│
├── 📂 docs/                        ← Documentação oficial
│   ├── 📂 product/                 ← MVP spec, plano de dev, perfil dev
│   └── 📂 tech/                    ← ⚠️ Arquivados (Unity) — ver ADR-004
│
├── 📂 server/                      ← 🆕 Servidor Go (a ser criado)
│   ├── main.go                     ← Entry point + tick loop
│   ├── world/                      ← Mundo, Blocos, NPCs
│   ├── combat/                     ← Motor D20, turnos, status
│   ├── economy/                    ← Inventário, crafting, comércio
│   └── data/                       ← JSONs e structs de dados
│
├── 📂 web/                         ← Cliente HTML/CSS/JS
│
├── 📂 vibe/                        ← Contexto evolutivo do projeto
│   ├── 📄 backlog.md               ← Tarefas por fase
│   ├── 📄 project-status.md        ← Snapshot do estado atual
│   ├── 📂 decisions/               ← ADRs (decisões de arquitetura)
│   │   ├── 📄 ADR-003-estrategia-repositorios.md  ← Ativo
│   │   └── 📄 ADR-004-pivot-mmorpg-servidor-go.md ← Ativo (PIVOT)
│   └── 📂 sessions/               ← Logs de sessão
│
└── 📂 pocs/                        ← ⚠️ Arquivo histórico (38 POCs Unity — substituídas pelas Fases)
```

---

## 🎯 Resumo Executivo

**O que fazer AGORA:**
1. **Fase 0 — Heartbeat:** `main.go` + tick + WebSocket

**O que fazer DEPOIS:**
2. **Fase 1 — Mundo Vivo:** NPCs com rotinas e IA
3. **Fase 2 — Observador:** Cliente web read-only
4. **Fase 3 — Jogador (MVP):** Criar, explorar, lutar, coletar

**O que NÃO fazer agora:**
- Tiers 2/3 das classes (Fase 5)
- Full Loot e Multiplayer (Fase 6)
- Livro de Magias (Fase 5)
- Grid tático (Fase 5)

---

> **Última atualização:** 2026-03-11
