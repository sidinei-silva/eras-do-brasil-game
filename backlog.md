# Backlog — Eras do Brasil (Game)

> Abra, veja a próxima tarefa, trabalhe nela.
> Detalhes completos (AC, referências GDD, folder structure): [docs/GUIA_RETOMADA.md](docs/GUIA_RETOMADA.md)
> Trilha de estudos Go e arquitetura: [estudos/README.md](estudos/README.md)

---

## Fase 0 — Heartbeat ✅

- [x] `go mod init` + estrutura de pastas (`server/`, `world/`, `engine/`, etc.)
- [x] `main.go` com game loop sequencial (`time.Ticker` + goroutine)
- [x] Struct `Mundo` em `world/mundo.go` com `ProcessarTick()`
- [x] WebSocket listener (`coder/websocket`) — cliente HTML recebe ticks via `world_snapshot`
- [x] Admin v0: `/admin/status` retorna status do game loop, uptime e online count

## Fase 1 — Mundo Vivo

- [ ] Structs NPC com rotinas diárias (acordar, trabalhar, comer, dormir)
- [ ] Utility AI básica (Score = Peso × (1 – Necessidade Normalizada))
- [ ] Ciclo dia/noite vinculado ao relógio do jogo (Manhã/Tarde/Noite/Madrugada)
- [ ] Sistema de Fofoca entre NPCs
- [ ] StoryManager — sementes narrativas reagem a threshold de eventos
- [ ] Admin v1: listar/localizar NPCs, inspecionar estado por entidade

## Fase 2 — Observador

- [ ] Cliente HTML/CSS/JS conecta via WebSocket, exibe estado do mundo
- [ ] Mapa de nós navegável (read-only, NPCs se movendo)
- [ ] Log de eventos em painel lateral
- [ ] HUD de tempo (relógio do jogo com períodos, dia/noite)
- [ ] Admin v2: observabilidade dos comandos no cliente

## Fase 3 — Jogador (MVP "O Despertar")

- [ ] Criação de personagem (1 classe: Guerreiro Tribal, point-buy 27 pts)
- [ ] Navegação por blocos com custo em tempo real (3 blocos: Vila, Floresta, Ruínas)
- [ ] Combate estático (D20 simplificado — Iniciativa → Turnos → Loot)
- [ ] Inventário e equipamentos (equipar, peso, capacidade)
- [ ] HUD principal (PV, XP, recursos, relógio, posição)
- [ ] Save/Load do estado do jogador (JSON)
- [ ] Admin v3: inspeção de personagem, inventário, save/load
- [ ] Playtest: loop completo 3+ vezes, ajustar números

## Fase 4 — Interação

- [ ] Diálogos ramificados com NPCs (árvore de diálogo JSON)
- [ ] Sistema de Quests (aceitar, rastrear, completar, recompensa)
- [ ] Crafting e coleta (proficiências, recursos, receitas)
- [ ] Comércio com NPCs (server-authoritative)
- [ ] Facções e reputação
- [ ] Status e condições em combate (envenenado, atordoado, queimando)
- [ ] Mini-campanha "O Caçador que Não Voltou" jogável completa
- [ ] Admin v4: facção, economia, diálogos, estado de quests

## Fase 5 — D20 Completo

- [ ] D20 completo (vantagem, desvantagem, críticos)
- [ ] Tiers 1→2→3 (Moedas de Classe, evolução, pré-requisitos)
- [ ] Herança de habilidades (Dom da Revivência)
- [ ] Habilidades ativas em combate (custo de recurso, AoE, recarga)
- [ ] 12 classes Tier 1 balanceadas
- [ ] Grid tático (posicionamento, AoE)
- [ ] Admin v5: simulação de combate e balanceamento

## Fase 6 — Multiplayer

- [ ] Múltiplas conexões WebSocket simultâneas
- [ ] Penalidade de morte (perda de 10% XP + 15% durabilidade — sem drop de itens)
- [ ] Quests competitivas (primeiro a entregar, timeout parcial, server-side)
- [ ] Inimigos evolutivos (Normal → Veterano → Alfa → Lenda, migram de região)
- [ ] Temporadas (state machine Tensão → Apogeu → Legado, votação simples)
- [ ] Eventos globais (mudanças de era)
- [ ] Economia multiplayer server-authoritative
- [ ] Bandeirantes como facção NPC ativa
- [ ] Admin v6: governança operacional + auditoria avançada

## Game Feel (Transversal)

- [ ] Loading/transições de tela (Fase 2+)
- [ ] Tutorial/onboarding em texto (Fase 3)
- [ ] Tela inicial + configurações web (Fase 3+)
- [ ] Sistema de áudio — Web Audio API (Fase 4+)
- [ ] Acessibilidade — contraste, fonte, teclado (Fase 5+)
- [ ] i18n — strings separadas do código (Fase 5+)

---

## Concluído

_Mover itens aqui quando prontos._
