# ADR-004: Pivot para MMORPG — Servidor Go + Cliente Web (MUD Moderno)

## Status

Aceito

## Contexto

O projeto Eras do Brasil passou por 3 anos de evolução documentária sem produzir código de jogo:

- **Jul 2023:** "1500 Caminhos do Brasil" — RPG single-player
- **2024 (Chats Gemini):** Visão expandida para MMORPG com servidor dedicado, tick global, NPCs estilo Dwarf Fortress, full loot, mundo que avança sem jogadores
- **Nov 2025 (commits 820ad6e, 16132ad, ed287af, c4fff4a):** Escopo reduzido para Co-op P2P com Unity — removeu StoryManager, missões competitivas, eventos globais, arquitetura servidor, full loot
- **Mar 2026:** 62 issues abertas, 38 POCs planejadas, 0 linhas de código de jogo. Developer travado por paralisia de documentação

**Diagnóstico:** Cada redução de escopo removeu as features que tornavam o jogo único e motivante para o criador. A desmotivação veio da perda de identidade do projeto, não do excesso de escopo.

**Perfil do developer:** Fullstack sênior (JS/TS/Go/PHP), 6h/semana, aprendendo Go.

## Decisão

**Pivotar para a visão original de MMORPG**, implementada como um **MUD Moderno** (servidor autoritário com cliente web baseado em texto):

### Stack

| Componente | Tecnologia | Justificativa |
|-----------|-----------|---------------|
| Servidor | **Go** | Goroutines para NPC AI concorrente, tick loops, múltiplas conexões. Sem bottleneck de event loop (Node.js) nem arquitetura request/response (PHP) |
| Cliente | **HTML/CSS/JS** (browser) | Zero atrito de instalação, usa skills existentes do developer |
| Protocolo | **WebSocket** | Comunicação bidirecional em tempo real |
| Dados estáticos | **Arquivos JSON** no repo (versionados) | Templates de NPC, items, mapa — carregados na RAM no startup |
| Dados dinâmicos | **SQLite** (WAL mode) via `modernc.org/sqlite` | RAM-first, snapshots async, zero ops — ver ADR-006 |

### Arquitetura

- **Servidor autoritário:** Toda lógica roda no servidor (ticks, NPC AI, combate, economia)
- **Tick global proativo:** `setInterval` dispara `Mundo.ProcessarTick()` independente de jogadores
- **Motor comutável:** Mesmo `ProcessarTick()` funciona offline (testes) e online (servidor)
- **Data-driven:** NPCs, items, quests definidos em JSON, motor genérico

### Fases de MVP

| Fase | Nome | Entrega |
|------|------|---------|
| 0 | Heartbeat | Servidor Go + tick loop + 1 NPC com rotina + log no terminal |
| 1 | Mundo Vivo | 5 NPCs, fofoca, dia/noite, eventos de tick |
| 2 | Observador | Cliente web mostrando estado do mundo via WebSocket |
| 3 | Jogador | Criar personagem, entrar no mundo, mover entre locais, ver NPCs |
| 4 | Interação | Falar com NPCs, aceitar missões, comprar itens, reputação |
| 5 | D20 | Combate por texto, rolagens D20, habilidades, morte/respawn |
| 6 | Multi | Segundo jogador conecta, mundo compartilhado |

### GDD

- **Rollback cirúrgico:** Restaurar seções MMO removidas nos 4 commits de Nov/Dez 2025
- **Adicionar do legado:** StoryManager, Missões Competitivas, BountyManager, Full Loot, Inimigos Evolutivos, Relógio da Ruptura
- **Preservar do atual:** Economia 5x5, Lifeskills/Proficiência, NPC AI (Utility AI), Facções 5 tiers, Magia (taxonomia completa), Combate (9 condições), Tick Debt (adaptado)
- **Formato:** Manter GDD narrativo como referência de design, criar specs técnicas em `docs/tech/` quando necessário para implementação

## Consequências

**Positivas:**

- Developer trabalha na stack onde é sênior (web + Go)
- Começa pelo que mais motiva (mundo vivo, NPC AI)
- Cada fase produz algo executável (não é mais só documentação)
- MUD moderno = scope mínimo de arte (texto), máximo de profundidade mecânica
- Go aprendido fazendo o que ama, não com exercícios genéricos
- Visão original restaurada — o jogo volta a ser O jogo dele

**Negativas:**

- MMORPG é scope infinito se não for disciplinado (mitigado por fases incrementais)
- Go é linguagem nova para o developer (mitigado: sintaxe simples, forte tipagem, boa documentação)
- Sem engine gráfica = limitação visual (mitigado: MUD moderno é o formato certo para solo dev)
- Precisa de hospedagem para servidor (mitigado: VPS barata, pode rodar em casa no início)

## Alternativas Consideradas

- **Manter Unity Co-op P2P:** Rejeitado — 3 anos sem código, developer desmotivado, Unity não é a skill principal
- **Node.js no servidor:** Rejeitado — event loop single-thread é bottleneck para simular centenas de NPCs pensando a cada tick
- **PHP no servidor:** Rejeitado — arquitetura request/response não serve para daemon persistente com WebSocket
- **Spin-offs primeiro:** Rejeitado — procrastinação disfarçada de preparação, aprende Go fazendo o jogo real
- **Repo novo:** Rejeitado — GDD atual é melhor que o legado (Caps 4-9 superiores), 91 commits de polimento, só precisa rollback cirúrgico

## Referências

- Tag `pre-pivot-mmorpg`: snapshot do estado Co-op P2P antes do pivot
- Repo legado: `sidinei-silva/ErasDoBrasil-Historico-Legado` (visão MMO original, Chats Gemini)
- Comparativo detalhado: `vibe/scratchpad/comparativo-gdd-legado-vs-atual.md`
- Commits de remoção MMO: `c4fff4a`, `820ad6e`, `16132ad`, `ed287af`
- ADR-005: Arquitetura do Servidor — Monolito com Goroutines e EventBus
- ADR-006: Persistência — RAM-First com Snapshots em SQLite
