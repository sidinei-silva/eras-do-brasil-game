# Lobby WebSocket na Pratica

Este guia explica o fluxo do lobby implementado no servidor.

## Visao geral do fluxo

1. Cliente abre conexao em /ws.
2. Handler faz upgrade HTTP -> WebSocket.
3. Hub cria um Client e registra no lobby.
4. Client inicia readPump e writePump em paralelo.
5. Mensagens recebidas no readPump viram comandos.
6. Hub faz broadcast para clientes via fila send.
7. writePump envia os eventos para o browser.

## Componentes

### Hub
Arquivo: `server/socket/hub.go`

Responsabilidades:

- manter o conjunto de clientes conectados
- registrar e remover clientes
- distribuir eventos para todos
- encerrar conexoes no shutdown

Canais usados:

- `register`
- `unregister`
- `broadcast`

## Client
Arquivo: `server/socket/client.go`

Responsabilidades:

- ler mensagens vindas do browser (`readPump`)
- escrever mensagens para o browser (`writePump`)
- validar eventos de entrada

Eventos de entrada atuais:

- `set_name`
- `chat`

Eventos de saida atuais:

- `player_joined`
- `player_left`
- `player_renamed`
- `chat`
- `world_snapshot`
- `error`

## Integracao com game loop

Arquivo: `server/main.go`

- game loop atualiza o mundo
- ticker paralelo envia `world_snapshot` para o hub
- endpoint `/admin/status` expoe `lobby.online`

## Mapa mental de concorrencia

- Hub.Run e dono do mapa `clients`
- readPump envia comandos
- writePump envia dados de saida
- OnlineCount usa `atomic` para leitura fora do Hub.Run

## Teste manual rapido

1. Suba o servidor: `cd server && go run main.go`
2. Abra dois clientes websocket (navegador ou extensão).
3. Envie `{"type":"set_name","name":"SeuNome"}`.
4. Envie `{"type":"chat","body":"ola"}`.
5. Observe eventos de join/chat/snapshot.

## Proximos passos de aprendizado

1. Extrair tipos de eventos para um pacote protocol.
2. Adicionar fila de comandos para o game loop.
3. Trocar broadcast global por broadcast por sala/regiao.
4. Criar testes de unidade para parse de eventos e regras do hub.
