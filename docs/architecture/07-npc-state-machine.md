# 07 — State Machine do NPC

> **Quando consultar:** Quando for implementar o NPCManager (Fase 1), quando precisar entender como Utility AI funciona, ou quando quiser adicionar um novo comportamento a um NPC.

## Diagrama — Decisão por Utility AI

```mermaid
flowchart TD
    START["NPCManager.ProcessTick()<br/>Para cada NPC:"] --> NEEDS

    NEEDS["Atualizar necessidades<br/><i>Fome +1, Energia -1, Social -1</i>"]
    NEEDS --> SCORE

    SCORE["Calcular prioridade de cada meta<br/><i>Comer: fome=70 → prioridade 70</i><br/><i>Trabalhar: agenda=50 → prioridade 50</i><br/><i>Dormir: energia=30 → prioridade 30</i>"]
    SCORE --> PICK

    PICK{"Meta de maior prioridade?"}
    PICK -->|"Comer (70)"| EAT["Ir pra taverna, comer"]
    PICK -->|"Trabalhar (50)"| WORK["Ir pro local de trabalho"]
    PICK -->|"Dormir (30)"| SLEEP["Ir pra casa, dormir"]
    PICK -->|"Socializar (25)"| SOCIAL["Encontrar outro NPC, fofoca"]

    EAT --> UPDATE["Atualizar estado do NPC<br/><i>Posição, necessidade, humor</i>"]
    WORK --> UPDATE
    SLEEP --> UPDATE
    SOCIAL --> UPDATE

    style START fill:#EEEDFE,stroke:#534AB7
    style PICK fill:#FAEEDA,stroke:#854F0B
```

## Diagrama — Rotina diária do Ferreiro (exemplo)

```mermaid
stateDiagram-v2
    [*] --> Dormindo: Madrugada

    Dormindo --> IndoPraForja: Manhã começa<br/>(se energia > 30)
    Dormindo --> Dormindo: Manhã começa<br/>(se energia <= 30, continua dormindo)

    IndoPraForja --> Trabalhando: Chegou na forja

    Trabalhando --> IndoPraTaverna: Fome > 60
    Trabalhando --> Trabalhando: Fome <= 60, continua

    IndoPraTaverna --> Comendo: Chegou na taverna
    Comendo --> VoltandoPraForja: Fome < 20

    VoltandoPraForja --> Trabalhando: Chegou na forja

    Trabalhando --> IndoPraCasa: Noite começa
    IndoPraCasa --> Dormindo: Chegou em casa

    note right of Trabalhando
        Agenda diz: "Manhã-Tarde: Forja"
        Mas necessidade pode quebrar agenda
        (Fome alta → vai comer primeiro)
    end note
```

## Diagrama — Sistema de fofoca

```mermaid
sequenceDiagram
    participant F as Ferreiro
    participant G as Guarda
    participant J as Jogador

    Note over F: Ferreiro viu lobo na floresta (manhã)
    F->>F: knowledgeBase += {tipo: "monstro",<br/>id: "lobo_01", local: "floresta_norte",<br/>visto_em: "manhã dia 3"}

    Note over F,G: Tarde — ambos na taverna
    F->>G: Fofoca: "Vi um lobo na floresta"
    G->>G: knowledgeBase += {tipo: "monstro",<br/>id: "lobo_01", local: "floresta_norte",<br/>fonte: "ferreiro", visto_em: "tarde dia 3"}

    Note over G,J: Jogador fala com o Guarda
    J->>G: "Sabe de algum perigo por aí?"
    G->>J: "O Ferreiro me contou que viu<br/>um lobo na Floresta do Norte"

    Note over F: Dia 5 — memória expira
    F->>F: knowledgeBase.remove(lobo_01)<br/>(visto há mais de 2 dias)
```

## Explicação — Utility AI vs State Machine

O NPC usa **dois sistemas combinados**:

A **Agenda** é o que o NPC "deveria" estar fazendo baseado no período do dia. O Ferreiro deveria trabalhar de Manhã a Tarde. É a rotina planejada, tipo um horário de trabalho.

As **Necessidades** são estados internos que decaem com o tempo (fome sobe, energia desce, social desce). Quando uma necessidade fica alta o suficiente, ela supera a prioridade da agenda. O Ferreiro deveria estar trabalhando, mas está com fome — ele "quebra" a rotina e vai comer.

A **Utility AI** é o mecanismo de decisão: calcula a prioridade de cada meta possível (trabalhar, comer, dormir, socializar) e escolhe a de maior prioridade. Não é um `if/else` hardcoded — é um sistema de scoring que permite comportamento emergente. Um NPC com o traço "Gula" tem peso maior pra fome, então vai comer mais cedo que outros.

## Necessidades e decaimento

Cada necessidade é um número de 0 a 100. A cada tick, as necessidades mudam:

| Necessidade | Decaimento por tick | O que sobe | O que desce |
|-------------|--------------------|-----------|-----------| 
| Fome | +1 por tick | Tempo sem comer | Comer na taverna |
| Energia | -1 por tick quando trabalhando | Dormir | Trabalhar, andar |
| Social | -0.5 por tick | Encontrar outro NPC | Ficar sozinho |
| Segurança | Variável | Estar em área perigosa | Estar em área segura |

Os valores exatos são tuning de game design — o importante é que o sistema existe e os valores são configuráveis por NPC (via traços).

## Traços do NPC

Traços são modificadores passivos que mudam como as necessidades evoluem ou como o NPC prioriza metas:

| Traço | Efeito |
|-------|--------|
| Gula | Fome cresce 50% mais rápido |
| Preguiçoso | Energia cai 50% mais rápido quando trabalhando |
| Ranzinza | Social cai mais devagar (não sente falta de companhia) |
| Valente | Segurança não afeta decisões (ignora perigo) |

## KnowledgeBase (o que o NPC sabe)

Cada NPC mantém uma base de conhecimento pessoal. Ele só sabe o que viu ou ouviu. Tipos de conhecimento:

| Tipo | Exemplo | Expira em |
|------|---------|-----------|
| Recurso | "Vi minério na mina" | 2 dias do jogo |
| Monstro | "Vi lobo na floresta" | 2 dias do jogo |
| Rotina de NPC | "O guarda vai pra taverna à noite" | 5 dias do jogo |
| Evento importante | "Jogador matou o lobo alfa" | Permanente |

O sistema de fofoca funciona quando dois NPCs se encontram no mesmo bloco: eles trocam 1-2 itens de suas knowledgeBases. Isso significa que informação se espalha organicamente pela rede social dos NPCs, sem que o jogador precise presenciar o evento diretamente.

## Para o jogador

O jogador acessa essa rede de informação conversando com NPCs. O que o NPC revela depende de:

1. **Afinidade**: o NPC gosta de você? (tabela de afinidade do Cap. 8.6)
2. **knowledgeBase**: o NPC realmente sabe o que você perguntou?
3. **Fofoca**: o NPC ouviu de outro? Nesse caso, a informação pode ser menos precisa.

## Implementação incremental

Fase 1 do projeto foca em NPC com agenda básica. A evolução:

1. **Fase 1**: NPC segue agenda fixa por período (Manhã: forja, Noite: casa). Sem necessidades.
2. **Fase 1 avançada**: Adiciona necessidades (fome, energia). NPC pode quebrar agenda.
3. **Fase 2**: Adiciona knowledgeBase. NPC acumula informação sobre o que vê.
4. **Fase 2 avançada**: Adiciona fofoca. NPCs trocam informação ao se encontrar.
5. **Fase 4**: Jogador pode conversar com NPC e acessar a knowledgeBase.
