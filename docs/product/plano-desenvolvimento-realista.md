# 🗓️ Plano de Desenvolvimento Realista — Eras do Brasil

**Data:** 2026-03-11 (atualizado pós-pivot)
**Premissa:** 3 noites/semana × ~2h/noite = **~6h/semana** (~24h/mês)
**Stack:** Servidor Go + Cliente Web (WebSocket)
**Referência:** [ADR-004 — Pivot MMORPG Servidor Go](../../vibe/decisions/ADR-004-pivot-mmorpg-servidor-go.md)

---

## Filosofia: Entregas Tangíveis, Sempre

> O maior risco do projeto não é técnico — é motivacional. Com tempo limitado, cada sessão precisa gerar algo **visível** e **jogável**. A regra de ouro: **nunca passe mais de 2 semanas sem algo rodando na tela.**

---

## Calendário Macro (2026)

| Mês | Foco Principal | Entregável |
|-----|---------------|------------|
| **Mar** | Pivot + Rollback GDD + Fase 0 | GDD atualizado, servidor Go com tick + WebSocket |
| **Abr** | Fase 1: Living World | 3 blocos, 2 NPCs com rotina por ticks |
| **Mai** | Fase 2: Observer | Cliente web mostrando mundo vivo em tempo real |
| **Jun** | Fase 3: Player (parte 1) | Criação de personagem, movimentação, Motor D20 |
| **Jul** | Fase 3: Player (parte 2) | Combate estático, inventário, save/load |
| **Ago** | Fase 4: Interaction (parte 1) | NPCs com diálogos, fofoca, facções |
| **Set** | Fase 4: Interaction (parte 2) | Crafting, coleta, economia, Relógio da Ruptura |
| **Out** | Fase 4: Interaction (parte 3) | 3 classes, 8 blocos, 1 mini-campanha |
| **Nov** | Playtest + polish | Loop completo funcional, testes |
| **Dez** | Retrospectiva + planejamento 2027 | Decisão sobre Fase 5 (D20 Full) |

> ⚠️ Este calendário é **otimista**. Se escorregar, a regra é: **terminar a fase atual antes de começar a próxima**.

---

## Detalhamento por Fase

### Fase 0: Heartbeat (Março 2026) — AGORA

**~12h (2 semanas)**
- [ ] `go mod init` + projeto Go básico rodando
- [ ] Tick loop com `time.Ticker` (goroutine)
- [ ] WebSocket server (`gorilla/websocket`)
- [ ] Cliente HTML mínimo conectado, mostrando "Tick: N"

### Fase 1: Living World (Abril 2026)

**~24h (4 semanas)**
- [ ] Structs para Bloco, NPC, Mundo
- [ ] Grafo de 3 blocos (Vila, Floresta, Ruínas) com adjacência
- [ ] NPC com agenda por tick (Ferreiro: casa→forja→casa)
- [ ] Utility AI básica (necessidades: fome, cansaço)
- [ ] Ciclo Dia/Noite por ticks
- [ ] Log de eventos: "Tick 42: Ferreiro chegou na Forja"

### Fase 2: Observer (Maio 2026)

**~24h (4 semanas)**
- [ ] Cliente web com mapa de nós (HTML/CSS clicável)
- [ ] Feed de eventos em tempo real via WebSocket
- [ ] Inspeção de NPC (clicar → ver estado, necessidade, ação)
- [ ] Relógio visual (Tick + período)

### Fase 3: Player (Junho-Julho 2026)

**~48h (8 semanas)**
- [ ] Criação de personagem (Origem + Classe + Atributos point-buy)
- [ ] Motor D20 (rolar + mods vs CD, acerto/falha/crítico)
- [ ] Movimentação por blocos (gastar Ticks)
- [ ] Combate estático por turnos (Iniciativa → D20 → Dano → Loot)
- [ ] Inventário (equipar/desequipar, peso)
- [ ] Save/Load (JSON)
- [ ] 1 Quest hardcoded ("O Caçador que Não Voltou")

### Fase 4: Interaction (Agosto-Outubro 2026)

**~72h (12 semanas)**
- [ ] Relógio da Ruptura (500 Ticks)
- [ ] NPCs com diálogos ramificados (condições)
- [ ] Fofoca (knowledgeBase entre NPCs)
- [ ] Facções e reputação
- [ ] Coleta + Crafting + Proficiências
- [ ] Comércio com NPCs
- [ ] Matriz de Itens 5×5
- [ ] 3 Origens × 1 Classe = 3 classes jogáveis
- [ ] Dom da Revivência (troca + herança de ativa)
- [ ] 8 blocos do Ato 1 navegáveis
- [ ] 1 mini-campanha completa

---

## Regras de Ouro para Sessões de Trabalho

### Estrutura de uma Noite (~2h)

```
[0:00-0:10] Setup: ler project-status.md, relembrar contexto
[0:10-1:45] Trabalho focado (1 tarefa por vez)
[1:45-2:00] Fechamento: atualizar status, anotar próximo passo
```

### Princípios

1. **Uma tarefa por sessão.** Começou, termina. Não deixe coisas pela metade.
2. **Sempre termine com algo rodando.** Pode ser um teste, pode ser um print, mas algo que prove que funciona.
3. **Use Vibe Keeper ao fim de cada sessão.** O `@vibe-keeper` vai salvar seu contexto para a próxima retomada.
4. **Se travou, mude para outra entrega da mesma fase.** Melhor fazer progresso lateral do que ficar parado.
5. **Não perfecione — itere.** Fase 0 não precisa ser perfeita. Precisa funcionar.

### Padrão de Commits

```
docs: resolve furo de lore — Dom da Revivência latente
gdd: atualiza Cap 8 para modelo servidor persistente
server: implementa tick loop + WebSocket
web: cria cliente mínimo com mapa de nós
vibe: sessão 2026-03-11 — pivot MMORPG rollback
```

---

## Métricas de Progresso

| Métrica | Como Medir | Meta Mensal |
|---------|-----------|-------------|
| POCs completas | Issues fechadas no GitHub | 2-3 POCs/mês |
| Sessões realizadas | Logs em vibe/sessions/ | 10-12/mês |
| Código novo | Linhas de C# no `game/` | Crescendo |
| Bugs/Furos resolvidos | Issues de Sprint 0 fechadas | 2-3/mês |
| Horas investidas | Estimativa por sessão | ~24h/mês |

---

*Este plano será revisado mensalmente pelo `@vibe-keeper` e ajustado conforme progresso real.*
