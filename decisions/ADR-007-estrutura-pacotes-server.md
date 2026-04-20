# ADR-007: Estrutura de Pacotes do Servidor — Flat por Domínio, Managers Donos do Estado

## Status

Aceito (2026-04-18). Revisado e consolidado (2026-04-20) após execução
da refatoração — ajustado para refletir a decisão final de **eliminar
o `state.GameState` central** em favor de managers donos do próprio
estado, com agregação via pacote `snapshot/`.

## Contexto

Após a sessão de reset (2026-04-03), o servidor foi reimplementado do
zero adotando uma estrutura DDD/Clean Architecture
(`internal/application`, `internal/domain`, `internal/infrastructure`).
A escolha foi pedagógica: permitir focar em aprender Go usando um
padrão arquitetural já familiar de sistemas CRUD enterprise.

Com o NPCManager implementado e a base do game loop estável, esse
andaime cumpriu seu propósito. Manter a estrutura DDD a partir dali
traria custos sem retorno proporcional:

1. **Idiomatismo Go conflita com DDD clássico.** Pacotes Go encapsulam
   por visibilidade, não por camada. Separar `application/npc` de
   `domain/npc` força tipos a serem públicos quando deveriam ser
   privados de pacote.
2. **Game servers têm acoplamento natural entre lógica e estado.**
   Separar isso em camadas técnicas é fricção sem benefício — a
   indústria (RuneScape emuladores, projetos Go de game server) usa
   layout flat.
3. **Imports verbosos.** `.../internal/application/world` adiciona
   ruído em todo arquivo novo.
4. **`internal/` e `cmd/` não justificados.** O servidor é um único
   binário standalone — semântica de `internal/` (impedir import
   externo) e `cmd/` (múltiplos entry points) não se aplicam.
5. **Contradição com ADRs vigentes.** ADR-005 e
   `docs/architecture/01-package-map.md` já especificavam layout flat
   por domínio — a estrutura DDD contradizia esses documentos.

Durante a execução da refatoração, uma **segunda decisão** precisou
ser tomada: um ciclo de dependência surgiu entre o `GameState` central
(que precisaria importar `npc.Npc`) e o `npc.Manager` (que precisaria
importar `state.GameState`). Esse ciclo é estrutural quando se tem um
god object central em Go — a solução não é quebrar com interfaces
(anti-idiomático quando só existe um tipo concreto), e sim repensar
a premissa: o `GameState` central não é necessário.

## Decisão

### Parte 1 — Layout flat por domínio

Adotar layout flat por domínio, com pacotes nascendo na raiz de
`server/`, agrupados por responsabilidade do jogo (não por camada
técnica).

### Parte 2 — Managers donos do próprio estado (sem god object)

Cada manager de domínio é dono do seu próprio estado. Não existe um
pacote `state/` central que agrega tudo. As consequências:

- `world.Manager` é dono do `GameTime`
- `npc.Manager` é dono do mapa de NPCs
- Futuros managers (`mob/`, `combat/`, ...) serão donos do seu estado
- O estado mutável é encapsulado — acesso externo é via métodos
  (`worldMgr.GameTime()`, `npcMgr.All()`, `npcMgr.Get(id)`)

Para agregar o estado de todos os managers num ponto (snapshot para
admin, persist), um pacote dedicado `snapshot/` faz esse papel:

- `snapshot.Build(tickCount, worldMgr, npcMgr, ...)` retorna struct
  imutável com cópia dos dados
- **É o único pacote que importa todos os managers**
- Ninguém importa `snapshot/` exceto `main.go` e consumidores
  externos (admin hub, persist)

### Estrutura final

```
server/
├── main.go              # entry point, wiring, graceful shutdown
├── engine/              # GameLoop (ticker + contador + reações)
├── command/             # PlayerCommand, CommandQueue, routers
├── world/               # Manager (dono do GameTime)
├── npc/                 # Manager (dono dos NPCs)
├── snapshot/            # Agregador transversal (Build)
├── net/                 # HTTP + WebSocket (player + admin)
├── data/                # loaders de JSON
└── persist/             # SQLite (Fase 1+)
```

Pacotes futuros (`mob/`, `combat/`, `story/`, `economy/`, `event/`)
nascem na raiz quando forem implementados.

### Regras de dependência

1. **`main.go` é o único que importa todos os managers** — faz
   wiring e registra reação do tick.
2. **Managers são donos do próprio estado.** Sem god object.
3. **Managers não importam comportamento de outros managers.** Podem
   importar **tipos universais** (ex: `npc` importa `world.GameTime`
   pra assinatura do `ProcessTick`). Acoplamento permitido é
   dado-a-dado, não lógica-a-lógica.
4. **`snapshot/` é o agregador.** Importa todos os managers. Ninguém
   importa `snapshot/` exceto `main.go` e consumidores externos de
   snapshot (admin hub, persist).
5. **`net/` e `persist/` não conhecem lógica de jogo.** Recebem
   `PlayerCommand` (via `command.Queue`) e entregam/recebem
   `Snapshot`.
6. **Pacote nasce quando tem habitante.** Não criar pacote vazio
   "por precaução".

### Por que sem `GameState` central (decisão revisada)

A versão original deste ADR propunha um pacote `state/` com um struct
`GameState` mutável central. Durante a implementação, esse desenho
gerou ciclo de dependência:

```
state  → npc   (GameState.NPCs precisa do tipo Npc)
npc    → state (Manager.ProcessTick(*GameState))
```

Três opções foram consideradas:

- **Pacote de tipos separado (`npctypes/`):** rejeitado — reintroduz
  separação dado/comportamento típica de DDD. Vira padrão pra todo
  domínio (`mobtypes/`, `combattypes/`). Métodos ficam num pacote
  diferente do struct — confuso.
- **Interface no `state/`:** rejeitado — polimorfismo onde só existe
  1 tipo concreto é overengineering. Em Go, interfaces vivem onde
  são consumidas, não onde os tipos vivem.
- **Sem `GameState`, managers donos do próprio estado:** adotado —
  elimina o ciclo sem sacrificar idioma Go, e ainda fornece
  encapsulamento real (mutações do estado passam por métodos do
  manager, não por mutação direta de map exposto).

O tradeoff é que o snapshot fica menos trivial — precisa orquestrar
chamadas a cada manager. Mas esse custo é pequeno e fica isolado no
`snapshot/Build()`.

## Consequências

**Positivas:**

- Imports curtos: `eras-do-brasil-game/server/world` em vez de
  `.../internal/application/world`
- Estrutura idiomática Go reconhecível por outros devs Go
- Sem ciclos de dependência estruturais
- Encapsulamento real — cada manager valida invariantes ao mutar seu
  estado
- Snapshot isolado num pacote que não contamina os demais
- Alinhamento código-documentação

**Negativas:**

- Adicionar um manager novo exige mexer em 2 lugares: `main.go`
  (wiring) e `snapshot/Build()` (agregação). Esperado — é o preço
  de não ter god object.
- Dev precisa resistir à tentação de colocar estado compartilhado
  em lugar central. Quando surgir algo que "todo mundo precisa" (ex:
  contador de tickets de combate), a pergunta é "qual manager é dono
  disso?", não "onde eu ponho esse global?".

## Plano de Execução (concluído em 2026-04-20)

Refatoração feita incrementalmente, cada passo um commit que compilou
e rodou:

1. ✅ Mover `internal/domain/world/` + `internal/application/world/`
   → `world/` (com `Manager` dono do `GameTime`)
2. ✅ Mover `internal/domain/npc/` + `internal/application/npc/`
   → `npc/` (com `Manager` dono dos NPCs, bootstrap via `LoadNpcsFromFile`)
3. ✅ Mover `internal/domain/command/` + `internal/application/command/`
   → `command/` (PlayerCommand, CommandQueue, routers, ProcessPlayerCommands)
4. ✅ Mover `internal/infrastructure/engine/` → `engine/` (só GameLoop)
5. ✅ Mover `internal/infrastructure/net/` → `net/`
6. ✅ Mover `internal/infrastructure/data/` → `data/`
7. ✅ Criar `snapshot/` com `Build()` agregando world + npc
8. ✅ Mover `cmd/game/main.go` → `main.go` (raiz de `server/`)
9. ✅ Eliminar `internal/domain/state/GameState` — estado redistribuído
   para os managers
10. ✅ Remover diretórios `internal/` e `cmd/`

Critério de "concluído" — todos atendidos:

- Pasta `internal/` não existe ✅
- Pasta `cmd/` não existe ✅
- Não existe pacote `state/` ✅
- `main.go` está na raiz de `server/` ✅
- `go build ./...` limpo ✅
- Servidor sobe, heartbeat funciona, player/admin conectam, tick
  avança, NPC atualiza necessidade e atividade ✅

## Alternativas Consideradas

- **Manter DDD/Clean Architecture:** Rejeitado — fricção crescente
  sem benefício correspondente para o caso de uso (game server solo).
- **Voltar a 1 arquivo só (`main.go` monolítico):** Rejeitado — a
  decisão da sessão de reset (2026-04-03) era válida no início, mas
  o número de domínios já justifica separação.
- **Usar `internal/` mas sem camadas DDD:** Rejeitado — `internal/`
  só faz sentido pra distinguir API pública de privada em bibliotecas.
  Servidor binário não tem essa preocupação.
- **Pacote `state/` central (versão original deste ADR):** Rejeitado
  após experiência prática — introduziu ciclo de dependência quando
  `npc.Manager` precisou importar `state` e `state` já importava `npc`.
- **`npctypes/`, `mobtypes/`, etc. para quebrar ciclo:** Rejeitado —
  ressuscita separação dado/comportamento típica de DDD. Anti-idiomático
  em Go e não generaliza bem.
- **Interfaces no `state/` para evitar import de tipos concretos:**
  Rejeitado — polimorfismo sem múltiplos tipos concretos é
  overengineering. Em Go, interfaces são definidas por quem consome,
  não por quem implementa.

## Referências

- ADR-005: Game loop sequencial (estrutura macro mantida)
- `docs/architecture/01-package-map.md`: diagrama de dependências atualizado
- Sessão 2026-04-03: contexto do reset que motivou a estrutura DDD
- Sessão 2026-04-18: análise inicial que levou à refatoração
- Sessão 2026-04-20: refatoração executada + decisão de eliminar `state/`
