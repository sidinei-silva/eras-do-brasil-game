# Plano: Schema e Dados dos Blocos (Feature 1.2)

## Context

A feature 1.2 (Blocos e Zonas) exige carregar o grafo estático do mundo via JSON no boot do servidor. Hoje existe um vácuo estrutural: o GDD descreve blocos **narrativamente** em vários arquivos (`05_Exploracao_e_Mundo.md`, `05B_Arquitetura_do_Mapa.md`, Atlas Ato 1), mas o `01_Schemas_Estruturais.md` **não tem schema formal para `Block**` — só tem `CharacterClass`, `Item`, `Inimigo`. O usuário pediu para:

1. **Revisar** se a estrutura atual dos blocos no GDD está boa.
2. **Criar** dados padronizados para 3 blocos, consumíveis pela feature 1.2, sem hardcode de valores no código Go.
3. **Definir** um padrão de ID e convenção de keys (em inglês).

## Review da estrutura atual (resposta à pergunta 1)

**Estado:**

- Schema **não existe** formalmente — é o principal gap.
- Campos necessários já aparecem **espalhados** no GDD: `type`, `travelCost`, `connections`, `services`, `rules`, `mobTable`, `resources`, `proceduralVariance`.
- Convenção de ID em uso no GDD para zonas: **snake_case sem prefixo** (`vila_sao_tome`, `floresta_norte`, `pico_da_neblina`). Consistente com o refinement 1.2. Divergente de NPCs (que usam `npc_*`), mas o contexto (arquivo `blocks.json`, struct `Block`) já deixa o tipo claro — o prefixo seria ruído.

**Veredito:** a informação narrativa é rica e consistente; falta só a **formalização**. É hora de criar o schema canônico.

## Decisões recomendadas

### 2.1 — Padrão de ID

- **Formato:** `snake_case` sem prefixo.
- **Exemplos:** `vila_sao_tome`, `floresta_do_norte`, `ruinas_queimadas`.
- **Justificativa:** o refinement 1.2 já adotou, o GDD já usa, e o contexto do arquivo/struct elimina a necessidade de prefixo `block_`.

### 2.2 — Convenção de keys

- **JSON:** `camelCase` (segue o padrão de `CharacterClass`, `Item`, `Inimigo` no GDD e dos DTOs de NPC no game).
- **Tudo em inglês**, incluindo valores de enum (`type: "urban"` e não `"Urbano"`) — facilita referência em código Go.

### 2.3 — Localização dos artefatos

- **Schema canônico (documentação):** nova seção `5.0 Schema: Block` em `gdd/06_Dados_e_Assets/01_Schemas_Estruturais.md` (repo `eras-do-brasil-gdd`).
- **Dados (implementação):** `server/data/blocks.json` (repo `eras-do-brasil-game`), seguindo o padrão do `npcs.json`.

### 2.4 — Escopo dos campos (balanço leve × completo)

O usuário pediu "todos os dados relevantes" mas "não muito pesado". Estratégia: schema com campos **obrigatórios** (consumidos pela 1.2) + **opcionais** (previstos para futuras features, mas já preenchíveis agora). Campos opcionais vazios = `[]` ou omitidos. O struct Go da 1.2 mapeia só o que consome; campos extras são ignorados pelo `encoding/json` sem erro — não pesa no build.

## Schema proposto (`Block`)

```json
{
  "id": "string",                   // snake_case, ex: "vila_sao_tome" — OBRIGATÓRIO
  "name": "string",                 // nome exibido, ex: "Vila de São Tomé" — OBRIGATÓRIO
  "region": "string",               // região pai, ex: "mata_costeira" — OBRIGATÓRIO
  "type": "string",                 // enum: "urban" | "dense_forest" | "mountain" | "water" | "ruins" | "wasteland" — OBRIGATÓRIO
  "levelRange": [0, 0],             // [min, max] nível recomendado — OBRIGATÓRIO
  "description": "string",          // descrição narrativa curta — OBRIGATÓRIO
  "connections": [                  // arestas saindo deste bloco — OBRIGATÓRIO
    {
      "toBlockId": "string",
      "travelMinutes": 10,          // custo em minutos de jogo (consumido pela 1.2)
      "terrain": "string"           // opcional, ex: "forest_path"
    }
  ],
  "services": [                     // serviços de hub, vazio se não aplicável — opcional
    {
      "id": "string",               // ex: "blacksmith"
      "name": "string",
      "actions": ["string"]
    }
  ],
  "rules": [                        // regras mecânicas ativas no bloco — opcional
    {
      "id": "string",
      "type": "string",             // ex: "improved_rest", "natural_cover"
      "effect": "string"
    }
  ],
  "mobTable": [                     // spawns, vazio em hubs — opcional
    {
      "enemyId": "string",
      "spawnChance": 0.5,
      "quantity": "1d2"
    }
  ],
  "resources": [                    // recursos coletáveis — opcional
    {
      "resourceId": "string",
      "rarity": "string",           // "common" | "uncommon" | "rare" | "epic" | "legendary"
      "spawnChance": 0.6,
      "quantity": "1d3"
    }
  ],
  "tags": ["string"]                // categorização livre, ex: ["safe", "hub"] — opcional
}
```

**Arquivo final:** `server/data/blocks.json` com wrapper `{ "blocks": [...] }`.

## Os 3 blocos (baseados no refinement 1.2 e Atlas do Ato 1)

1. `**vila_sao_tome**` — Hub seguro, nível 0, tipo `urban`. Tem `services` (blacksmith, tavern, chapel) e `rules` (improved_rest). Conexões: → `floresta_do_norte` (10min).
2. `**floresta_do_norte**` — Exploração intermediária, nível 1–3, tipo `dense_forest`. Tem `rules` (natural_cover), `mobTable` (lesser_spirit, wolf), `resources` (iron_wood). Conexões: → `vila_sao_tome` (10min), → `ruinas_queimadas` (15min).
3. `**ruinas_queimadas**` — Zona final do grafo, nível 3–5, tipo `ruins`. Tem `rules` (spiritual_echoes), `mobTable` (bandeirante_scout), `resources` (spiritual_ashes). Conexões: → `floresta_do_norte` (15min).

**Topologia:** caminho linear `vila_sao_tome ↔ floresta_do_norte ↔ ruinas_queimadas`, conexões bidirecionais e simétricas (travelMinutes iguais nos dois sentidos), satisfazendo todos os invariantes da seção 7 do refinement 1.2.

## Arquivos críticos

### Repo `eras-do-brasil-gdd` (branch `claude/review-gdd-files-A4dwq`)

- **Editar:** `gdd/06_Dados_e_Assets/01_Schemas_Estruturais.md` — adicionar seção `5.0 Schema: Block` com JSON schema + tabela de análise de campos (mesmo formato das seções 2/3/4 existentes).

### Repo `eras-do-brasil-game` (branch `claude/review-gdd-files-A4dwq`)

- **Criar:** `server/data/blocks.json` — contendo os 3 blocos completos.

**Não serão tocados nesta entrega:**

- Código Go (`server/world/manager.go`, DTOs, loader) — faz parte da implementação da feature 1.2, fora do escopo deste pedido (o usuário pediu **dados**, não implementação).

## Verificação (end-to-end)

1. **Schema no GDD:** visualmente, confirmar que a nova seção `5.0 Schema: Block` segue o mesmo formato das seções `2.0`, `3.0`, `4.0` (JSON schema + tabela "Análise dos Campos").
2. **JSON válido:** `cat server/data/blocks.json | jq .` deve retornar sem erro.
3. **IDs consistentes:** os 3 IDs (`vila_sao_tome`, `floresta_do_norte`, `ruinas_queimadas`) batem exatamente com os listados no refinement 1.2 (seção 4 e 6).
4. **Simetria de conexões:** para cada par `(A, B, mins)`, existe também `(B, A, mins)` com o mesmo `travelMinutes`.
5. **Sem self-connection:** nenhum bloco tem `toBlockId == id`.
6. **Futura validação via 1.2:** quando a feature 1.2 for implementada, o loader deve conseguir desserializar o JSON sem erro e `worldMgr.GetBlock("floresta_do_norte")` deve retornar o bloco com 2 conexões.

## Referências

- **Padrão de schema a seguir:** `gdd/06_Dados_e_Assets/01_Schemas_Estruturais.md` — seções 2.0/3.0/4.0 (CharacterClass, Item, Inimigo).
- **Contrato mínimo:** `docs/product/refinements/1.2-blocos-zonas.md` — seções 4 (escopo), 6 (critérios de aceite), 7 (invariantes).
- **Conteúdo narrativo dos 3 blocos:** `gdd/05_Livros_Auxiliares/01_Atlas_do_Eco_Ato1.md` e `gdd/03_Enredo_e_Mundo/01_Ato_1_A_Primeira_Ruptura.md`.
- **Padrão de dados no game:** pasta `server/data/` e formato do `npcs.json` (tags `json:` em camelCase, carga única no boot).
