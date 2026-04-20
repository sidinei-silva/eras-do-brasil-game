# ADR-007: Estrutura de Pacotes do Servidor — Flat por Domínio

## Status

Aceito (2026-04-18)

## Contexto

Após a sessão de reset (2026-04-03), o servidor foi reimplementado do zero
adotando uma estrutura DDD/Clean Architecture (`internal/application`,
`internal/domain`, `internal/infrastructure`). A escolha foi pedagógica:
permitir ao desenvolvedor focar em aprender Go usando um padrão
arquitetural já familiar de sistemas CRUD enterprise.

Com o NPCManager implementado e a base do game loop estável, esse
andaime cumpriu seu propósito. Manter a estrutura DDD a partir daqui traz
custos sem retorno proporcional:

1. **Idiomatismo Go conflita com DDD clássico.** Pacotes Go encapsulam por
   visibilidade, não por camada. Separar `application/npc` de `domain/npc`
   força tipos a serem públicos quando deveriam ser privados de pacote.
2. **Game servers têm acoplamento natural entre lógica e estado.** Separar
   isso em camadas técnicas é fricção sem benefício — a indústria
   (RuneScape emuladores, projetos Go de game server) usa layout flat.
3. **Imports verbosos.** `.../internal/application/world` adiciona ruído
   em todo arquivo novo.
4. **`internal/` e `cmd/` não justificados.** O servidor é um único binário
   standalone — semântica de `internal/` (impedir import externo) e
   `cmd/` (múltiplos entry points) não se aplicam.
5. **Contradição com ADRs vigentes.** ADR-005 e
   `docs/architecture/01-package-map.md` já especificam layout flat por
   domínio. A estrutura DDD atual contradiz esses documentos.

## Decisão

Adotar **layout flat por domínio**, com pacotes nascendo na raiz de
`server/`, agrupados por responsabilidade do jogo (não por camada técnica).

### Estrutura

```
server/
├── main.go              # entry point, wiring, graceful shutdown
├── engine/              # GameLoop, CommandQueue
├── state/               # GameState — estado mutável central
├── command/             # tipos de comando + roteadores
├── world/               # WorldManager + GameTime + Block
├── npc/                 # NPCManager + NPC + Agenda + Need
├── net/                 # HTTP + WebSocket (player + admin)
├── data/                # loaders de JSON e templates estáticos
└── persist/             # SQLite (Fase 1+)
```

Pacotes futuros (`mob/`, `combat/`, `story/`, `economy/`, `event/`)
nascem na raiz quando forem implementados.

### Regras de dependência (mantidas do ADR-005 e package-map)

1. `main.go` é o único arquivo que importa todos os pacotes.
2. Managers de domínio (`world`, `npc`, ...) **nunca** importam outro
   manager. Trocam informação via `state.GameState`.
3. Pacotes de I/O (`net`, `persist`) não conhecem lógica de jogo.
4. `state/` e `command/` não importam managers — só são importados.
5. Pacote nasce quando tem habitante. Não criar pacote vazio
   "por precaução".

### Por que `state/` separado em vez de `shared/`

A versão original do package-map propunha um pacote `shared/` agrupando
`GameState`, `Command`, `Event`. Foi descartado porque:

- `state.GameState` é a **entidade central mutável** do servidor — merece
  pacote próprio pela centralidade conceitual, não pra "agrupar tipos".
- `command/` tem comportamento próprio (Router, Queue) e vai crescer.
  Não é só tipo.
- `event/` ainda não tem habitante. Cria quando o EventBus voltar (se
  voltar) para notificações async.

## Consequências

**Positivas:**

- Imports curtos: `eras-do-brasil-game/server/world` em vez de
  `.../internal/application/world`
- Estrutura idiomática Go reconhecível por outros devs Go
- Alinhamento código-documentação (ADR-005, package-map)
- Concluído o ciclo pedagógico de aprendizado iniciado no reset

**Negativas:**

- Refatoração custa ~3-5 sessões de trabalho de movimentação de arquivos
  e ajuste de imports
- Tipos que eram públicos entre `application/` e `domain/` precisam ser
  reavaliados quanto à visibilidade (trabalho de limpeza posterior)

## Plano de Execução

Refatoração incremental, cada passo um commit que compila e roda:

1. Mover `internal/domain/state/` → `state/`
2. Mover `internal/domain/command/` → `command/` (mesclar com router)
3. Mover `internal/domain/world/` + `internal/application/world/` → `world/`
4. Mover `internal/domain/npc/` + `internal/application/npc/` → `npc/`
5. Mover `internal/infrastructure/engine/` → `engine/`
6. Mover `internal/infrastructure/net/` → `net/`
7. Mover `internal/infrastructure/data/` → `data/`
8. Mover `cmd/game/main.go` → `main.go` (raiz de `server/`)
9. Remover diretórios vazios (`internal/`, `cmd/`)

Critério de "concluído":

- Pasta `internal/` não existe
- Pasta `cmd/` não existe
- `main.go` está na raiz de `server/`
- `go build ./...` limpo
- `go vet ./...` limpo
- Servidor sobe, heartbeat funciona, player/admin conectam, tick avança,
  NPC atualiza necessidade

## Alternativas Consideradas

- **Manter DDD/Clean Architecture:** Rejeitado — fricção crescente sem
  benefício correspondente para o caso de uso (game server solo).
- **Voltar a 1 arquivo só (`main.go` monolítico):** Rejeitado — a
  decisão da sessão de reset (2026-04-03) era válida no início, mas
  o número de domínios já justifica separação.
- **Usar `internal/` mas sem camadas DDD:** Rejeitado — `internal/`
  só faz sentido para distinguir API pública de privada em bibliotecas.
  Servidor binário não tem essa preocupação.

## Referências

- ADR-005: Game loop sequencial (estrutura macro mantida)
- `docs/architecture/01-package-map.md`: regras de dependência
- Sessão 2026-04-03: contexto do reset que motivou a estrutura DDD
- Sessão 2026-04-18: análise que levou a esta decisão
