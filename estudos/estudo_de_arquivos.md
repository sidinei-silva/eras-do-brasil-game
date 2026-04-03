# Guia de estudo dos arquivos do servidor

Objetivo: estudar o codigo ja escrito do mais simples para o mais complexo, de dentro para fora, registrando seus comentarios com sua propria assinatura.

## Como usar este guia

1. Abra um arquivo por vez, na ordem abaixo.
2. Leia e responda as perguntas sugeridas para cada etapa.
3. Escreva seus comentarios no seu caderno de estudo (ou abaixo de cada item, se quiser manter tudo aqui).
4. Marque o item como concluido quando terminar.

## Modelo de comentario pessoal

Use este formato para manter consistencia nas suas anotacoes:

- Nome: Seu Nome
- Data: AAAA-MM-DD
- Arquivo: caminho/do/arquivo.go
- O que eu entendi:
- Duvidas que ficaram:
- Termos de Go para revisar:
- Proximo passo:

## Trilha de estudo (ordem recomendada)

- [ ] [server/npc/schedule.go](../server/npc/schedule.go) - apenas um type alias, 4 linhas
- [ ] [server/npc/npc.go](../server/npc/npc.go) - NPC, Needs, Activity, sem deps externas
- [ ] [server/npc/registry.go](../server/npc/registry.go) - colecao de NPCs, map em Go
- [ ] [server/npc/loader.go](../server/npc/loader.go) - como carregar JSON embutido
- [ ] [server/world/zone.go](../server/world/zone.go) - structs de zona
- [ ] [server/world/map.go](../server/world/map.go) - grafo de zonas
- [ ] [server/world/world.go](../server/world/world.go) - o mundo, mutex, tick
- [ ] [server/engine/publisher.go](../server/engine/publisher.go) - interface SnapshotPublisher
- [ ] [server/engine/gameloop.go](../server/engine/gameloop.go) - o loop, atomic, context
- [ ] [server/socket/socket.go](../server/socket/socket.go) - apenas um handler
- [ ] [server/socket/client.go](../server/socket/client.go) - uma conexao WS
- [ ] [server/socket/hub.go](../server/socket/hub.go) - channels, select, observer
- [ ] [server/admin/commands.go](../server/admin/commands.go) - roteamento de comandos
- [ ] [server/admin/hub.go](../server/admin/hub.go) - hub admin
- [ ] [server/main.go](../server/main.go) - a cola de tudo

## Perguntas-guia por etapa

Use estas perguntas em cada arquivo:

1. Quais tipos e responsabilidades este arquivo define?
2. Que dependencias internas ele usa (pacotes do projeto)?
3. Onde entram goroutines, canais, mutex ou context?
4. Que dado entra e que dado sai deste arquivo?
5. Que parte ficou confusa para revisar depois?

## Depois dos comentarios: mao no volante

Este bloco veio de uma orientacao de IA para te ajudar a retomar o projeto com autonomia. A ideia aqui e usar como roteiro pratico, escrevendo o codigo voce mesmo.

### 1) Leitura ativa com papel e caneta (ou comentarios no codigo)

Ordem sugerida para este ciclo:

- [server/npc/npc.go](../server/npc/npc.go)
- [server/npc/schedule.go](../server/npc/schedule.go)
- [server/npc/registry.go](../server/npc/registry.go)
- [server/npc/loader.go](../server/npc/loader.go)
- [server/world/world.go](../server/world/world.go)

Para cada funcao, escreva com suas palavras o que ela faz.

### 2) Teste de compreensao com mudancas pequenas (feito por voce)

- Renomear loop para gameLoop (mudanca controlada e didatica).
- Adicionar um novo campo no Snapshot (exemplo: quantidade de NPCs).
- Adicionar endpoint HTTP /admin/npcs retornando mundo.NPCs.AllStates().

### 3) Proximo passo real do backlog

A Utility AI basica e o item mais educativo agora. Formula base:

Score = Peso x (1 - NecessidadeNormalizada)

Voce vai precisar entender bem os Needs do NPC (Hunger, Fatigue) para implementar. Essa e uma parte importante para destravar o restante da Fase 1.

### 4) Antes de qualquer nova feature

Adicione primeiro o endpoint /admin/npcs. Ele te obriga a conectar mentalmente:

- Registry -> World
- Rota HTTP -> leitura do estado atual

## Proximo passo concreto sugerido

Implemente voce mesmo, sem IA gerando codigo:

GET /admin/npcs retornando o estado atual de todos os NPCs em JSON.

Checklist tecnico para implementar:

- Entender [server/npc/registry.go](../server/npc/registry.go) e como funciona AllStates().
- Criar http.HandleFunc em [server/main.go](../server/main.go).
- Serializar resposta com json.NewEncoder(w).Encode(...), no mesmo padrão dos endpoints existentes.

Observação: sao poucas linhas, mas cada linha toca um conceito importante do projeto.