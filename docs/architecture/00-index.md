# Documentação de Arquitetura — Eras do Brasil (Game Server)

> Referência técnica para estudo e consulta durante o desenvolvimento.
> Cada documento contém diagramas Mermaid (renderizam no GitHub e VS Code) + explicação teórica.

## Como usar

Estes documentos servem dois propósitos: entender a arquitetura antes de codar e consultar durante o desenvolvimento quando surgir dúvida sobre "onde coloco isso" ou "como isso se conecta".

A ordem de leitura recomendada segue a numeração dos arquivos. Os primeiros (01-05) cobrem a visão macro — como os pacotes se organizam, como o tick funciona, como I/O se conecta à lógica. Os seguintes (06-10) detalham ciclos de vida específicos. Os últimos (11-13) cobrem infraestrutura e patterns de Go.

## Índice

| # | Arquivo | O que responde |
|---|---------|----------------|
| 01 | [Mapa de pacotes](01-package-map.md) | Onde coloco esse código? Quem importa quem? |
| 02 | [Anatomia de um tick](02-tick-anatomy.md) | O que acontece a cada 500ms? Em que ordem? |
| 03 | [Fluxo WebSocket](03-websocket-flow.md) | Como o cliente se comunica com o game loop? |
| 04 | [Combate dentro do tick](04-combat-lifecycle.md) | O combate é paralelo? Como ele avança? |
| 05 | [Manager pattern](05-manager-pattern.md) | Como organizar um pacote internamente? |
| 06 | [Ciclo de vida da conexão](06-connection-lifecycle.md) | Do upgrade HTTP ao disconnect — onde fica o defer? |
| 07 | [State machine do NPC](07-npc-state-machine.md) | Como o NPC decide o que fazer? |
| 08 | [Ciclo de vida do mob](08-mob-lifecycle.md) | Spawn, patrulha, aggro, morte, respawn |
| 09 | [Fluxo de dados](09-data-flow.md) | Estáticos vs dinâmicos, JSON vs RAM vs SQLite |
| 10 | [Jornada de um comando](10-player-command-flow.md) | Do clique no browser até a resposta na tela |
| 11 | [Mapa de goroutines](11-goroutine-map.md) | Todas as goroutines ativas, quem cria, quem mata |
| 12 | [Graceful shutdown](12-graceful-shutdown.md) | Sequência de desligamento seguro |
| 13 | [Patterns de Go no projeto](13-go-patterns.md) | defer, ctx, channels, select — quando usar cada um |

## Convenções

Os diagramas usam Mermaid e renderizam automaticamente no GitHub. No VS Code, instale a extensão "Markdown Preview Mermaid Support".

Cada documento segue a estrutura: diagrama → explicação do que ele mostra → regras/patterns extraídos → quando consultar este documento.
