# 01 — Mapa de Dependência de Pacotes

> **Quando consultar:** Quando surgir a dúvida "onde coloco esse código?" ou "esse pacote pode importar aquele?"

## Diagrama

```mermaid
graph TD
    subgraph "Ponto de entrada"
        MAIN["main.go<br/><i>startup, wiring, shutdown</i>"]
    end

    subgraph "Motor do jogo"
        ENGINE["engine/<br/><i>GameLoop, EventBus</i>"]
    end

    subgraph "Managers de domínio (lógica de jogo)"
        WORLD["world/<br/><i>Tempo, blocos, clima</i>"]
        NPC["npc/<br/><i>IA, rotinas, fofoca</i>"]
        MOB["mob/<br/><i>Spawn, patrulha, aggro</i>"]
        COMBAT["combat/<br/><i>D20, turnos, dano</i>"]
        STORY["story/<br/><i>Quests, eventos, arcos</i>"]
        ECONOMY["economy/<br/><i>Preços, crafting, trade</i>"]
    end

    subgraph "I/O (goroutines paralelas)"
        NET["net/<br/><i>HTTP, WebSocket, Hub</i>"]
        PERSIST["persist/<br/><i>SQLite snapshots</i>"]
    end

    subgraph "Tipos compartilhados"
        SHARED["shared/<br/><i>GameState, Command, Entity</i>"]
    end

    MAIN --> ENGINE
    MAIN --> WORLD
    MAIN --> NPC
    MAIN --> MOB
    MAIN --> COMBAT
    MAIN --> STORY
    MAIN --> ECONOMY
    MAIN --> NET
    MAIN --> PERSIST

    ENGINE -.-> SHARED
    WORLD -.-> SHARED
    NPC -.-> SHARED
    MOB -.-> SHARED
    COMBAT -.-> SHARED
    STORY -.-> SHARED
    ECONOMY -.-> SHARED
    NET -.-> SHARED
    PERSIST -.-> SHARED
```

## Explicação

O `main.go` é o único arquivo que conhece todos os pacotes. Ele cria cada manager, injeta dependências, e passa tudo pro game loop. Nenhum outro arquivo importa tantos pacotes.

Os **managers de domínio** (world, npc, mob, combat, story, economy) nunca se importam entre si. Se o NPCManager precisa saber que horas são, ele não faz `import "world"` — ele recebe essa informação via `GameState`, um struct definido no pacote `shared/` que todos conhecem. Isso evita dependências circulares e mantém cada pacote testável isoladamente.

Os **pacotes de I/O** (net, persist) também não conhecem a lógica de jogo. O pacote `net/` sabe receber JSON e colocar um `Command` numa fila. O pacote `persist/` sabe pegar entidades dirty e gravar no SQLite. Nenhum dos dois sabe o que é um NPC, um combate, ou uma quest.

O **pacote shared/** contém os tipos que todo mundo precisa: `GameState` (o snapshot do mundo que passa entre managers), `Command` (um comando do jogador parseado), `Entity` (interface base para qualquer coisa persistível), `Event` (notificação do EventBus). Este pacote não tem lógica — só definições de tipos.

## Regras para adicionar um pacote novo

Quando for criar um pacote novo, pergunte:

1. **Ele processa lógica de jogo por tick?** → Vai na camada de managers, recebe `ProcessTick()`, é chamado pelo game loop.
2. **Ele faz I/O (rede, disco, banco)?** → Vai na camada de I/O, roda em goroutine própria, comunica via channel.
3. **Ele define tipos que outros pacotes precisam?** → Vai no `shared/`.
4. **Ele precisa importar outro manager?** → Não deveria. Repense a responsabilidade ou passe a informação via `GameState`.

## Regra de ouro

A seta de dependência só aponta para baixo (do main pro resto) e para o shared (de qualquer um pro shared). Nunca para cima (manager importando main) e nunca horizontalmente (manager importando outro manager).
