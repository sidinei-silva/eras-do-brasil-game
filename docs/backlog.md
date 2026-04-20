# Backlog — Eras do Brasil

> Reimplementação do zero. Voce decide a ordem, a estrutura e o como.
> Quando travar, pesque do commit `4d64406`.

---

## Fase 0 — Heartbeat

**Objetivo:** Servidor Go rodando com game loop que avança o tempo do mundo. Clientes recebem estado via WebSocket.

- [x] Servidor HTTP escuta na porta 8080
- [x] Game loop roda em intervalo fixo, avançando o tempo do jogo
- [x] Mundo tem hora do jogo (comeca em 1500-01-01 06:00) e periodo (Manha/Tarde/Noite/Madrugada)
- [x] Cliente HTML conecta via WebSocket e recebe estado do mundo a cada tick
- [x] Endpoint `/admin/status` retorna status do servidor em JSON

## Fase 1 — Mundo Vivo

**Objetivo:** NPCs vivem no mundo com rotinas, necessidades e comportamento autonomo.

- [x] NPCs tem agenda por periodo (onde estar, o que fazer)
- [x] NPCs tem necessidades (fome, cansaco) que mudam a cada tick
- [ ] Utility AI decide atividade baseado nas necessidades (Score = Peso x (1 - Necessidade))
- [ ] NPCs fofocam entre si (propagacao de informacao)
- [ ] Admin pode listar e inspecionar NPCs

## Fase 2 — Observador

**Objetivo:** Cliente web mostra o mundo em tempo real (read-only).

- [ ] Mapa de nos navegavel mostrando NPCs se movendo
- [ ] Log de eventos lateral
- [ ] HUD de tempo (relogio do jogo, periodo)

## Fase 3 — Jogador (MVP "O Despertar")

**Objetivo:** Loop jogavel: criar personagem, explorar, lutar, completar quest, salvar.

- [ ] Criar personagem (Guerreiro Tribal, point-buy 27 pts)
- [ ] Navegar entre 3 zonas (Vila, Floresta, Ruinas)
- [ ] Combate estatico D20 simplificado
- [ ] Inventario e equipamentos
- [ ] Quest "O Cacador que Nao Voltou" jogavel
- [ ] Save/Load do estado

> Detalhes: [mvp-o-despertar-spec.md](product/mvp-o-despertar-spec.md)

## Fase 4+ — Futuro

Dialogos, crafting, faccoes, D20 completo, multiplayer. Detalhar quando chegar aqui.

---

## Referencias

- ADRs: [decisions/](../decisions/)
- GDD: [eras-do-brasil-gdd](https://github.com/sidinei-silva/eras-do-brasil-gdd)
- Go reference: [estudos/go/](../estudos/go/)
- Codigo anterior (backup): commit `4d64406`
