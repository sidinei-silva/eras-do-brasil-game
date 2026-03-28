# Referencia Rapida - Conceitos Go

## Ponteiros: * e &

### O que e um ponteiro
Um ponteiro guarda o endereco de memoria de um valor.

```go
valor := 42
endereco := &valor
```

### * (desreferenciar)

```go
var p *int = &x
*p = 100
fmt.Println(x)
```

### & (pegar endereco)

```go
x := 42
p := &x
```

### Por que usar ponteiros

1. Compartilhar estado entre funcoes e goroutines.
2. Evitar copia de structs grandes.
3. Manter identidade unica de uma entidade em memoria.

No projeto:

- `func NewWorld() *World`
- `func NewGameLoop(interval time.Duration, world *world.World) *GameLoop`
- `func NewClient(hub *Hub, conn *websocket.Conn) *Client`

## Concorrencia: goroutines, channels e select

### Goroutine

```go
go func() {
    fmt.Println("Roda em paralelo")
}()
```

No projeto:

- `go loop.Start(ctx)` em `server/main.go`
- `go hub.Run(ctx)` em `server/main.go`
- `go client.readPump()` e `go client.writePump()` em `server/socket/hub.go`

### Channel

```go
ch := make(chan string)
go func() { ch <- "oi" }()
msg := <-ch
```

No projeto:

- `register chan *Client`
- `unregister chan *Client`
- `broadcast chan []byte`
- `send chan []byte`

### Select

```go
select {
case <-ctx.Done():
    return
case msg := <-ch:
    fmt.Println(msg)
}
```

No projeto:

- Hub processa `register/unregister/broadcast` em `Run`.
- Game loop alterna entre `ctx.Done()` e `ticker.C`.

## Context

Context e o sinal de cancelamento entre goroutines.

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

No projeto:

- `signal.NotifyContext` inicia o contexto raiz em `server/main.go`.
- `Hub.Run(ctx)` encerra clientes no shutdown.
- `GameLoop.Start(ctx)` encerra o ticker no shutdown.

## Sincronizacao: mutex vs atomic

### Mutex
Use mutex quando varios campos precisam mudar juntos.

### Atomic
Use atomic para um valor simples e independente.

No projeto:

- `running atomic.Bool` e `tickCount atomic.Uint64` no game loop.
- `online atomic.Int64` no hub.
- `sync.RWMutex` no estado de mundo.

## WebSocket no modelo de lobby

### Regra de ouro
Uma goroutine escreve, outra goroutine le.

- `readPump` faz apenas leitura do socket.
- `writePump` faz apenas escrita do socket.

Isso evita corridas de escrita e estados inconsistentes da conexao.

### Heartbeat ping/pong

- Servidor manda ping periodico.
- Cliente responde pong automaticamente.
- Se pong nao chega no prazo, conexao e encerrada.

## Defer

`defer` agenda limpeza para o fim da funcao.

No projeto:

- fechamento de conexao
- stop de ticker
- cancel de contexto

## Erros comuns para evitar

1. Escrever no mesmo websocket por mais de uma goroutine.
2. Esquecer de fechar ticker.
3. Bloquear em envio de canal sem fallback.
4. Atualizar mapa compartilhado fora da goroutine dona do estado.

## Leituras recomendadas

- Effective Go: https://go.dev/doc/effective_go
- Go Memory Model: https://go.dev/ref/mem
- Go Concurrency Patterns: https://go.dev/blog/pipelines
