# Backlog — Eras do Brasil

> **Como usar este backlog:**
>
> - Cada item é uma linha curta. Detalhamento acontece **no momento em que você pega o item**, usando `docs/product/refinement-template.md`.
> - Pré-requisitos explicitados onde há dependência de features não-implementadas.
> - Quando um item tem refinamento feito, o link vai na linha.
> - Se travar, pesque do commit `4d64406`.

---

## Fase 0 — Heartbeat ✅

**Objetivo:** Servidor Go rodando com game loop que avança o tempo. Clientes recebem estado via WebSocket.

- [x] Servidor HTTP escuta na porta 8080
- [x] Game loop roda em intervalo fixo, avançando o tempo do jogo
- [x] Mundo tem hora do jogo (começa em 1500-01-01 06:00) e período (Manhã/Tarde/Noite/Madrugada)
- [x] Cliente HTML conecta via WebSocket e recebe estado do mundo a cada tick
- [x] Endpoint `/admin/status` retorna status do servidor em JSON
- [x] Refatoração de arquitetura: flat por domínio, managers donos do próprio estado (ADR-007)

## Fase 1 — Mundo Vivo

**Objetivo:** NPCs vivem no mundo com rotinas, necessidades e comportamento autônomo. Mundo persistente simulando vida mesmo sem jogador.

**Ordem de execução:** NPCs primeiro → Blocos/Zonas depois → Integração (movimento real entre blocos) por último.

### 1.1 — NPCs (base) ✅

- [x] NPCs têm agenda por período (onde estar, o que fazer)
- [x] NPCs têm necessidades (Fome, Energia) que mudam a cada tick
- [x] Decisão de atividade por threshold hardcoded (if Hunger >= 80, etc.)

### 1.2 — Blocos e Zonas **PRÓXIMO**

> **Refinamento:** [refinements/1.2-blocos-zonas.md](product/refinements/1.2-blocos-zonas.md)

- [x] Struct `Block` com ID, nome, tipo (Vila/Floresta/Ruínas), descrição
- [x] Grafo de conectividade entre blocos (adjacência, custo de travessia em minutos de jogo)
- [x] 3 blocos iniciais da Mata Costeira ligados (Vila de São Tomé → Floresta do Norte → Ruínas Queimadas)
- [x] `world.Manager` dono dos blocos; getters pro snapshot
- [x] Template JSON de bloco + loader
- [x] Admin pode listar blocos e ver conexões

### 1.3 — NPCs (Utility AI intermediário) — ✅

> **Refinamento:** [refinements/1.3-utility-ai-intermediario.md](product/refinements/1.3-utility-ai-intermediario.md)

- [x] Adicionar need `Social` (decaimento na ausência de outros NPCs no mesmo zone string)
- [x] Adicionar need `Seguranca` (placeholder fixo, sem efeito real até blocos existirem)
- [x] Trocar thresholds hardcoded por sistema de scoring (Score = Peso × f(need))
- [x] Pesos configuráveis via JSON no template do NPC (sem Traits ainda)
- [x] Função de decisão retorna a meta vencedora; fallback Idle se nenhuma passar do threshold mínimo
- [x] Admin pode ver scoring de um NPC específico (comando OOB retorna ranking de metas)

### 1.4 — Movimento de NPCs (POIs & validação de localização) — Andamento

> Refinamento detalhado: [refinamento 1.4 — POIs e Localização Validada de NPCs](docs/product/refinements/1.4-movimento-de-npcs.md)

- Escopo principal: introduzir POIs intra-bloco (lista de IDs por bloco) e validar localizações de NPCs no load; não implementa trânsito entre blocos nem estado walking.
- Entradas:
  - [ ] `Block` ganha campo `Pois []string` (loader + snapshot)
  - [ ] Refactor dos campos de NPC: `currentZone` → `currentBlock`, adicionar `currentPoi` (vazio = outdoor), `homePoi`, `eatingPoi`; `schedule[].location` → `schedule[].poiId`
  - [ ] Loader valida cruzamento `npcs.json` ↔ `blocks.json` (erros fatais na inicialização com mensagem clara)
  - [ ] `TransitionTo(activity, poiId, gameTime)` atualiza `CurrentPoi`; `CurrentBlock` NÃO muda nesta fase
  - [ ] Regra de Loneliness: comparar (CurrentBlock, CurrentPoi); `CurrentPoi==""` = "outdoor" virtual
  - [ ] Comandos admin: `admin_get_npc_full` retorna `currentBlock` + `currentPoi`; novo `admin_get_pois_in_block <blockId>`
  - [ ] Migração de dados: mapear valores antigos para `currentBlock`/`poiId` PT-BR snake_case (ex.: `oficina_do_ferreiro`, `taverna_da_vila`)

- Exceções (o que NÃO entra nesta tarefa):
  - walking / trânsito intra-bloco com tempo
  - movimento entre blocos (itinerantes)
  - POI como struct rica (descrição/capacity/ownership)
  - persistência do POI atual entre reinícios
- Critérios de aceite resumidos:
  - [ ] Servidor sobe limpo com `npcs.json` e `blocks.json` migrados
  - [ ] Servidor **falha na inicialização** se `currentBlock` ou qualquer `poiId` referenciar entidade inexistente, mostrando NPC e campo inválido
  - [ ] Exemplos funcionais: Tomas mantém `currentBlock = vila_sao_tome` e transita entre POIs do bloco conforme schedule; regras de Loneliness seguem D3 (POIs diferentes = separados; ambos `""` no mesmo bloco = juntos)
  - [ ] `admin_get_npc_full` e `admin_get_pois_in_block` retornam campos esperados
  - [ ] `go build` e `go vet` limpos

### 1.5 — Fofoca (requer 1.4 completo)

- [ ] `knowledgeBase` do NPC: lista de `Knowledge{tipo, id, local, visto_em}`
- [ ] Evento "viu X" popula knowledgeBase quando NPC chega num bloco com recurso/mob/outro NPC
- [ ] Quando 2+ NPCs no mesmo bloco por N ticks: trocam 1-2 itens de knowledgeBase
- [ ] Esquecimento: entradas expiram após prazo por tipo (recurso 2 dias, rotina 5 dias, etc.)
- [ ] Admin pode consultar knowledgeBase de um NPC

### 1.6 — Ferramentas de dev/admin

- [ ] Listar NPCs via comando admin (já parcialmente implementado)
- [ ] Inspecionar NPC individual (stats, needs, atividade, zona, knowledgeBase)
- [ ] Forçar avanço de N períodos (pular dia) pra testar rotinas rápido
- [ ] Dump do mundo em JSON (tick, período, todos NPCs, todos blocos)

---

## Fase 2 — Observador

**Objetivo:** Cliente web mostra o mundo em tempo real (read-only). Jogador "assiste" o mundo viver.

- [ ] Mapa de nós navegável mostrando NPCs se movendo entre blocos
- [ ] Log de eventos lateral (NPC chegou em X, viu Y, fofocou com Z)
- [ ] HUD de tempo (relógio do jogo, período, dia do jogo)
- [ ] Painel de inspeção de NPC (clicar mostra needs, atividade, knowledgeBase)

## Fase 3 — Jogador (MVP "O Despertar")

**Objetivo:** Loop jogável: criar personagem, explorar, lutar, completar quest, salvar.

- [ ] Criar personagem (Guerreiro Tribal, point-buy 27 pts)
- [ ] Jogador entra em bloco → vê NPCs e recursos do bloco
- [ ] Combate estático D20 simplificado (ataque, defesa, HP)
- [ ] Inventário e equipamentos (estrutura mínima)
- [ ] Quest "O Caçador que Não Voltou" jogável end-to-end
- [ ] Save/Load do estado (primeiro uso do `persist/` + SQLite)

> Detalhes: [mvp-o-despertar-spec.md](product/mvp-o-despertar-spec.md)

## Fase 4+ — Futuro

Diálogos, crafting, facções, D20 completo, multiplayer, Traits dos NPCs, 5° need (Sede), economia de recursos. Detalhar quando chegar aqui.

---

## Referências

- ADRs: [decisions/](../decisions/)
- GDD: [eras-do-brasil-gdd](https://github.com/sidinei-silva/eras-do-brasil-gdd)
- Go reference: [estudos/go/](../estudos/go/)
- Código anterior (backup): commit `4d64406`
- Template de refinamento: [refinement-template.md](product/refinement-template.md)
