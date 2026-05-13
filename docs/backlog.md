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

### 1.4 — Movimento de NPCs (POIs & validação de localização) — ✅

> Refinamento detalhado: [refinamento 1.4 — POIs e Localização Validada de NPCs](docs/product/refinements/1.4-movimento-de-npcs.md)

- Escopo principal: introduzir POIs intra-bloco (lista de IDs por bloco) e validar localizações de NPCs no load; não implementa trânsito entre blocos nem estado walking.
- Entradas:
  - [x] `Block` ganha campo `Pois []string` (loader + snapshot)
  - [x] Refactor dos campos de NPC: `currentZone` → `currentBlock`, adicionar `currentPoi` (vazio = outdoor), `homePoi`, `eatingPoi`; `schedule[].location` → `schedule[].poiId`
  - [x] Loader valida cruzamento `npcs.json` ↔ `blocks.json` (erros fatais na inicialização com mensagem clara)
  - [x] `TransitionTo(activity, poiId, gameTime)` atualiza `CurrentPoi`; `CurrentBlock` NÃO muda nesta fase
  - [x] Regra de Loneliness: comparar (CurrentBlock, CurrentPoi); `CurrentPoi==""` = "outdoor" virtual
  - [x] Comandos admin: `admin_get_npc_full` retorna `currentBlock` + `currentPoi`; novo `admin_get_pois_in_block <blockId>`
  - [x] Migração de dados: mapear valores antigos para `currentBlock`/`poiId` PT-BR snake_case (ex.: `oficina_do_ferreiro`, `taverna_da_vila`)

- Exceções (o que NÃO entra nesta tarefa):
  - walking / trânsito intra-bloco com tempo
  - movimento entre blocos (itinerantes)
  - POI como struct rica (descrição/capacity/ownership)
  - persistência do POI atual entre reinícios
- Critérios de aceite resumidos:
  - [x] Servidor sobe limpo com `npcs.json` e `blocks.json` migrados
  - [x] Servidor **falha na inicialização** se `currentBlock` ou qualquer `poiId` referenciar entidade inexistente, mostrando NPC e campo inválido
  - [x] Exemplos funcionais: Tomas mantém `currentBlock = vila_sao_tome` e transita entre POIs do bloco conforme schedule; regras de Loneliness seguem D3 (POIs diferentes = separados; ambos `""` no mesmo bloco = juntos)
  - [x] `admin_get_npc_full` e `admin_get_pois_in_block` retornam campos esperados
  - [x] `go build` e `go vet` limpos

### 1.5 — Fofoca / KnowledgeBase (requer 1.4 completo) **PRÓXIMO**

> **Refinamento:** [refinements/1.5-fofoca-knowledge-base.md](product/refinements/1.5-fofoca-knowledge-base.md)

- Escopo: substrato funcional de KnowledgeBase com fofoca entre os 3 residentes da Vila. Schema único enriquecido (`SeenCount + FirstSeenAt + LastSeenAt + LearnedAt + Source + Important`). Auto-geração só para `AVISTAMENTO_NPC` (único tipo automático nesta fase).
- Versão honesta do escopo original — itinerantes, StoryManager e outros tipos de Knowledge ficam adiados conscientemente para 1.6+.
- Entradas:
  - [ ] Struct `Knowledge` com 10 campos; enum `KnowledgeType` com único valor `AVISTAMENTO_NPC`
  - [ ] `Npc` ganha `KnowledgeBase []Knowledge` + `LastGossipedWith map[string]time.Time`
  - [ ] `knowledge_config.json` com expiração per-tipo, cooldown de fofoca, tamanho máximo, quantidade trocada por evento
  - [ ] Auto-geração: quando 2+ NPCs no mesmo `(Block, Poi)`, cada um popula entry sobre o outro; dedup por chave `(Type, EntityId, BlockId, PoiId)` incrementa `SeenCount` ao invés de duplicar
  - [ ] Fofoca: ≥2 NPCs co-localizados + cooldown por par expirado → troca 1-2 entries aleatórias; anti-loop via filtro de `Source` e `EntityId`
  - [ ] Expiração: lazy nos getters + sweep no início do `ProcessTick`; entries `Important=true` isentas
  - [ ] FIFO de tamanho: descarta entry com `LearnedAt` mais antigo se exceder limite; entries `Important` isentas
  - [ ] Snapshot inclui cópia profunda da KB
  - [ ] Comandos admin: `admin_inject_knowledge` (strict validation), `admin_get_npc_knowledge`, `admin_get_npc_full` ganha `knowledgeBaseSize`
  - [ ] Flag de log `NPCKnowledge`
- Exceções (o que NÃO entra — débitos conscientes):
  - Tipos `LOCAL`, `RECURSO`, `ENTIDADE`, `EVENTO_IMPORTANTE` auto-gerados (só `AVISTAMENTO_NPC`; outros entram com seus respectivos geradores em 1.6+ / Fase 2+)
  - Itinerantes / NPCs inter-bloco (1.6)
  - StoryManager (Fase 2 ou 3)
  - Diálogo do jogador / acesso da KB via NPC (Fase 3+)
  - Sistema de Afinidade / Relacionamentos NPC↔NPC (Fase 3+)
  - Inferência de rotina (`ROTINA_NPC` como tipo derivado) (Fase 3+ ou conforme demanda)
  - Compartilhamento de `SeenCount` via fofoca (decisão 3.9 do refinamento)
  - Persistência da KB entre reinícios
  - Traits que modificam expiração / decay (Fase 4+)
  - Separação visto/conhecimento como duas estruturas (refatorável quando demanda surgir)
- Critérios de aceite resumidos:
  - [ ] Servidor sobe com `knowledge_config.json` válido; falha com mensagem clara se inválido
  - [ ] Após 1 dia de jogo: Tomas e Naila têm entries um sobre o outro (no mínimo via almoço/jantar na taverna); re-encontros incrementam `SeenCount` sem duplicar
  - [ ] Fofoca propaga entries respeitando cooldown e anti-loop (`Source != receptor`)
  - [ ] Sweep remove expiradas por `LearnedAt`; `Important=true` sobrevive
  - [ ] Admin: `admin_inject_knowledge` strict, `admin_get_npc_knowledge` retorna entries ativas
  - [ ] `go build` e `go vet` limpos

### 1.6 — Itinerantes (NPCs inter-bloco)

> **Refinamento:** _a redigir antes de iniciar_

- Escopo: introduzir NPCs que mudam de `CurrentBlock` ao longo do dia (caçadores, mercadores, guardas em patrulha). Refactor para separar `HomeBlock` de `CurrentBlock`. Estado de trânsito (`walking`) entre blocos com tempo configurável. Propagação natural de KB entre vilas via itinerantes que carregam fofoca de um lado pro outro.
- Entradas (rascunho — refinar):
  - [ ] `Npc` ganha `HomeBlock` separado de `CurrentBlock`
  - [ ] Activity nova: `ActivityTravelingBetweenBlocks` com tempo configurável
  - [ ] Schedule pode incluir POIs de blocos diferentes (validação cruzada amplia o que a 1.4 faz)
  - [ ] Pelo menos 1 NPC itinerante novo nos `npcs.json` (ex: caçador que sai pra floresta de manhã, volta à noite)
  - [ ] Loneliness e fofoca já funcionam (1.5) — o ganho é que itinerantes naturalmente propagam KB entre blocos
  - [ ] Persistência de POI atual? (decidir no refinamento — provavelmente ainda não)
- Exceções:
  - Pathfinding intra-bloco (mapa local, Camada 3 — T2+)
  - Travel cancelado por evento (entrará junto com StoryManager)
- Critérios de aceite resumidos:
  - [ ] Caçador novo sai da Vila pela manhã, chega na Floresta do Norte após N min de jogo, volta à tarde
  - [ ] KB do caçador acumula entries em ambos os blocos
  - [ ] Encontro do caçador com Tomas/Ricardo/Naila na Vila propaga conhecimento da Floresta
  - [ ] `go build` e `go vet` limpos

### 1.7 — Ferramentas de dev/admin

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
