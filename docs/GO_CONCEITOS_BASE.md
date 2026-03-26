# 🧠 Referência Rápida — Conceitos Go Usados na Fase 0

## Ponteiros: `*` e `&`

### O que é um ponteiro?
Um **ponteiro** é um endereço de memória — o "local" onde um valor está guardado.

```
valor := 42          // um número 42 guardado em algum lugar da memória
endereco := &valor   // & pega o endereço de valor (algo como 0x7fff5fbff8c8)
```

### `*` (asterisco) = "acesse o valor apontado"

```go
var p *int = &x    // p é um ponteiro para int
*p = 100           // mude o valor no endereço para 100
fmt.Println(x)     // imprime 100 (x foi modificado!)
```

### `&` (ampersand) = "pegue o endereço"

```go
x := 42
p := &x           // p aponta para x
```

### Por que usar ponteiros?
1. **Compartilhar estado:** se você passa uma cópia, mudanças não afetam original
2. **Eficiência:** não copia dados grandes (e.g., um struct de 10MB)
3. **Ciclo de vida:** permitir que múltiplas goroutines trabalhem no mesmo objeto

Em nosso código:
- `func NewWorld() *World` — retorna um PONTEIRO, não uma cópia
- `loop.world *world.World` — loop tem acesso ao MESMO mundo, mudanças são vistas por todos

---

## Concorrência: Goroutines, Channels, Select

### Goroutine = "thread leve" de Go

```go
go func() {
    fmt.Println("Roda em paralelo!")
}()
// código continua aqui sem esperar a goroutine terminar
```

Em nosso código:
- `go loop.Start(ctx)` — game loop roda em paralelo com main
- `go func() { server.ListenAndServe() }()` — servidor HTTP roda em paralelo

### Context = "sinal de cancele"

Um context é um mecanismo para sinalizar que algo deve parar.

```go
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(1 * time.Second)
    cancel()  // cancela o context
}()

<-ctx.Done()  // bloqueia até ctx ser cancelado
fmt.Println("Cancelado!")
```

Em nosso código:
- O context vem do sinal do SO (Ctrl+C)
- `<-ctx.Done()` bloqueia main até o SO enviar sinal
- `loop.Start(ctx)` vê `<-ctx.Done()` e para o loop

### Channel = "tubo de comunicação" entre goroutines

```go
ch := make(chan string)
go func() {
    ch <- "oi"  // envia no channel
}()
msg := <-ch     // recebe do channel (bloqueia até algo chegar)
fmt.Println(msg)  // imprime "oi"
```

Em nosso código:
- `ticker.C` é um channel que envia a cada 500ms
- `<-ticker.C` bloqueia até o próximo tick

### Select = "espere vários channels, responda o primeiro que chegar"

```go
select {
case <-ctx.Done():
    fmt.Println("Cancelado")
    return
case <-ticker.C:
    fmt.Println("Tick!")
}
```

Em nosso código:
- Select espera DOIS canais:
  1. `ctx.Done()` — se SO cancelar
  2. `ticker.C` — se chegar nova hora de tick
- Qualquer um que "dispare" primeiro, aquele case executa

---

## Sincronização: Mutex e Atomic

### Problema: múltiplas goroutines modificando o mesmo dado

```go
var count = 0

go func() { count++ }()
go func() { count++ }()

// count pode ser 0, 1, ou 2?
// Depende de qual goroutine terminou primeira — é indeterminado!
```

### Solução 1: sync.Mutex (Mutual Exclusion)

```go
var mu sync.Mutex
var count = 0

go func() {
    mu.Lock()    // trava
    count++
    mu.Unlock()  // destrava
}()
```

**Regras do mutex:**
- `Lock()` e `Unlock()` SEMPRE em pares
- Não use `Lock()` depois `RUnlock()` — erro!
- Se esquecer `Unlock()`, o programa trava para sempre

Em nosso código:
```go
g.mu.Lock()
g.lastTickDuration = time.Since(start)
g.mu.Unlock()
```

### Solução 2: sync.RWMutex (Read-Write Mutex)

Permite **múltiplos leitores** OU **um escritor** (não ambos).

```go
var mu sync.RWMutex

// Múltiplas goroutines podem fazer isto ao mesmo tempo:
go func() {
    mu.RLock()
    fmt.Println(value)
    mu.RUnlock()
}()

// Mas só uma goroutine pode fazer isto:
go func() {
    mu.Lock()
    value = 42
    mu.Unlock()
}()
```

Em nosso código:
```go
g.mu.RLock()      // leitor bloqueia
last := g.lastTickDuration
g.mu.RUnlock()    // libera
```

### Solução 3: sync/atomic — para uma variável simples

```go
var running atomic.Bool

// Seguro para múltiplas goroutines:
running.Store(false)
is := running.Load()
```

**Quando usar:**
- `atomic.Bool`, `atomic.Uint64` — para um valor simples
- `sync.Mutex` — para vários valores que sempre modificam juntos
- `sync.RWMutex` — quando há muitos leitores e poucos escritores

Em nosso código:
- `running atomic.Bool` — só precisa de true/false, nada complexo
- `tickCount atomic.Uint64` — só precisa de um contador
- `mu sync.RWMutex` para `lastTickDuration` — pode ter muitos Status() lendo enquanto atualiza

---

## Padrões de Código

### Defer = "execute no final"

```go
func abre_arquivo() {
    f := open("dados.txt")
    defer f.Close()  // garante que Close() é chamado ao sair
    
    // faz coisas com f
    // mesmo se houver erro ou return precoce, Close() é chamado
}
```

Em nosso código:
- `defer stop()` — garante que cancela o context
- `defer ticker.Stop()` — garante que desliga o ticker
- `defer cancel()` — garante que cancela o shutdown context

### Interface = "contrato"

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Qualquer tipo que tiver método Read() é um Reader
```

Na Fase 1 vamos usar interfaces para managers (NPC, Combat, etc.) terem todos `ProcessTick()`.

---

## Teste Rápido de Compreensão

**Pergunta 1:** Se mudo `go loop.Start(ctx)` para `loop.Start(ctx)` (sem `go`), o que muda?
- Resposta: main() esperaria que loop.Start() terminasse (nunca — loop infinito), então nunca chegaria em `<-ctx.Done()`.

**Pergunta 2:** Se tirar `defer g.mu.Unlock()`, o que acontece?
- Resposta: na primeira leitura de Status(), o mutex fica travado para sempre — deadlock.

**Pergunta 3:** Por que `w.mu.Lock()` e depois `defer w.mu.Unlock()` em ProcessTick?
- Resposta: garante que ninguém lê enquanto atualiza tick e gameTime.

---

## Links de Estudo

- **Go Concurrency Patterns:** https://go.dev/blog/pipelines
- **Effective Go:** https://go.dev/doc/effective_go
- **The Go Memory Model:** https://go.dev/ref/mem
