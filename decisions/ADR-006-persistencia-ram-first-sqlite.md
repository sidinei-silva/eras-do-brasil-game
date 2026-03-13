# ADR-006: Persistência — RAM-First com Snapshots em SQLite

## Status

Aceito

## Contexto

Um servidor de jogo tem requisitos de I/O fundamentalmente diferentes de uma aplicação web:

- **Web app:** request → query DB → response. Latência de DB é aceitável (cada request é independente).
- **Game server:** o mundo processa ticks a cada ~500ms com centenas de entidades. Qualquer I/O no game loop travaria o tick.

Precisamos de uma estratégia que:

1. Mantenha o game loop em **zero I/O** (tudo na RAM)
2. Persista estado de forma **assíncrona** (sem bloquear gameplay)
3. Seja **simples de operar** (solo dev, 6h/semana)
4. Permita **migração futura** se o projeto crescer

## Decisão

### Princípio: RAM é a fonte de verdade durante gameplay. O banco é um checkpoint.

### Dois Tipos de Dados

| Tipo | Exemplos | Armazenamento | Quando Carrega |
|------|----------|---------------|----------------|
| **Estático** (read-only) | Templates de NPC, items, skills, mapa, quests | Arquivos JSON no repositório (versionados via git) | Startup — carrega tudo na RAM |
| **Dinâmico** (read-write) | Posição do jogador, inventário, estado do mundo, preços, memória NPC | RAM durante gameplay, SQLite para persistir | Startup — carrega do SQLite pra RAM |

### Banco de Dados: SQLite (WAL mode)

**Por que SQLite:**

- **Zero ops** — sem Docker, sem connection pool, sem processo separado
- **Embarcado** — compila junto com o binário Go via `modernc.org/sqlite` (Go puro, sem CGO)
- **WAL mode** — leituras e escritas concorrentes sem lock global
- **1 arquivo** — backup = `cp game.db game.db.bak`
- **SQL é SQL** — migrar pra PostgreSQL depois é trocar o driver e ajustar queries

**Lib Go:** `modernc.org/sqlite` — implementação Go pura, sem dependência de C/CGO, cross-compile sem problemas.

### Padrão de Escrita (Zero I/O no Game Loop)

```
[Goroutines do Jogo]
     │
     │  modificam structs na RAM
     │  marcam entidade como "dirty"
     │
     ▼
[Channel de Persistência]   ← buffered channel
     │
     ▼
[PersistManager goroutine]
     │
     │  coleta dirty entities
     │  batch INSERT/UPDATE em transação
     │  limpa dirty flags
     │
     ▼
[SQLite WAL mode]
```

**Snapshot periódico:** A cada 30-60 segundos, PersistManager coleta todas as entidades dirty e faz batch write dentro de uma transação SQLite.

**Escrita crítica:** Ações irreversíveis (trade entre jogadores, morte com full loot, level-up) enviam imediatamente para o channel de persistência. A escrita continua sendo async — nunca bloqueia o tick — mas entra na fila com prioridade.

**Shutdown gracioso:**
1. `signal.Notify(SIGTERM, SIGINT)` captura sinal
2. Para de aceitar novas conexões
3. Flush de todas as entidades dirty para SQLite
4. Fecha conexão SQLite
5. `os.Exit(0)`

### Padrão de Leitura

- **Na startup:** `SELECT * FROM players`, `SELECT * FROM world_state`, `SELECT * FROM npcs` → popula structs Go na RAM
- **Durante gameplay:** ZERO queries ao banco. Toda leitura vem da RAM.

### Schema Inicial

```sql
-- Jogadores
CREATE TABLE players (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    block_id TEXT NOT NULL,
    data JSON NOT NULL,
    updated_at INTEGER NOT NULL
);

-- Estado global do mundo
CREATE TABLE world_state (
    key TEXT PRIMARY KEY,
    data JSON NOT NULL,
    updated_at INTEGER NOT NULL
);

-- NPCs com estado dinâmico
CREATE TABLE npcs (
    id TEXT PRIMARY KEY,
    block_id TEXT NOT NULL,
    data JSON NOT NULL,
    updated_at INTEGER NOT NULL
);

-- Índices para queries de startup
CREATE INDEX idx_players_block ON players(block_id);
CREATE INDEX idx_npcs_block ON npcs(block_id);
```

**Estratégia de colunas:** Campos que precisam de query SQL (block_id para saber quem está onde) ficam em colunas dedicadas. Dados complexos (inventário, stats, memória NPC) ficam em colunas JSON — flexíveis, sem migrations para cada campo novo.

### O que NÃO fazer

- **Persistir em arquivos JSON** — escrita concorrente corrompe, sem transação atômica, reescreve arquivo inteiro para mudar 1 campo
- **Query ao banco durante gameplay** — latência de disco mata o tick loop
- **PostgreSQL agora** — overhead operacional desnecessário para 1 developer com <100 jogadores

### Caminho de Migração para PostgreSQL

Quando (e **se**) precisar:

1. >100 jogadores simultâneos E queries analíticas em tempo real
2. Trocar `modernc.org/sqlite` por `pgx`
3. Ajustar sintaxe SQL (mínima — ambos suportam JSON)
4. PersistManager continua idêntico — só muda o driver

SQLite em WAL mode aguenta centenas de escritas por segundo. Para Fases 0-6 do MVP, é mais que suficiente.

## Consequências

**Positivas:**

- Game loop com zero I/O — performance máxima para ticks
- SQLite embarcado — zero infraestrutura, backup trivial
- Schema flexível com JSON — evolui sem migrations pesadas
- Snapshot async — persistência sem impacto no gameplay
- Shutdown gracioso — zero perda de dados em deploy/restart
- Dados estáticos versionados no git — rollback = `git checkout`

**Negativas:**

- Crash sem graceful shutdown perde dados do último intervalo de snapshot (mitigado: intervalo curto de 30-60s, ações críticas entram na fila imediatamente)
- Coluna JSON não tem schema enforcement no banco (mitigado: structs Go tipadas são o schema, validação acontece no código)
- SQLite não escala para múltiplas instâncias de servidor (mitigado: irrelevante — decisão de escalar é para quando houver need real, ADR futuro)

## Alternativas Consideradas

- **PostgreSQL desde o início:** Rejeitado — overhead operacional (Docker, connection pooling, migrations, backups programados) desproporcional para solo dev em fase MVP. Pode ser adotado no futuro sem reescrever a lógica de negócio.
- **Arquivos JSON para estado dinâmico:** Rejeitado — sem transações atômicas, sem escrita concorrente segura, sem partial update. Aceitável apenas para dados estáticos (templates).
- **Redis como cache intermediário:** Rejeitado — adiciona dependência externa e processo separado. RAM do próprio processo Go já é o "cache" — Redis seria redundante.
- **BoltDB / BadgerDB (key-value Go):** Rejeitado — sem SQL, sem JSON queries, mais difícil de inspecionar manualmente. SQLite permite `sqlite3 game.db` para debug.

## Referências

- ADR-004: Pivot para MMORPG — define stack e fases do MVP
- ADR-005: Arquitetura do Servidor — define o PersistManager como goroutine
- GDD Cap 8 §8.12: Arquitetura do Motor de Mundo
