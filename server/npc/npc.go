package npc

// Activity descreve o que um NPC está fazendo num dado momento.
// Usando tipo nomeado sobre string (igual ao ZoneType que você já viu)
// para o compilador rejeitar strings arbitrárias onde Activity é esperado.
type Activity string

const (
	ActivitySleep  Activity = "sleep"
	ActivityWork   Activity = "work"
	ActivityEat    Activity = "eat"
	ActivityPatrol Activity = "patrol"
	ActivityIdle   Activity = "idle"
)

// Needs representa as necessidades biológicas/sociais do NPC.
// Valores float64 de 0.0 (satisfeito) a 1.0 (crítico).
//
// Por que float64 e não int?
// Precisamos de frações pequenas por tick (ex: 0.00069) para que a fome
// leve ~1 dia de jogo para chegar ao máximo. Com int, o incremento mínimo
// seria 1, o que zeraria o NPC em 100 ticks (~50 segundos reais).
type Needs struct {
	Hunger  float64 // fome: sobe a cada tick, cai ao comer
	Fatigue float64 // cansaço: sobe ao trabalhar/andar, cai ao dormir
}

// decayRates define o quanto cada necessidade muda por tick.
// 1 tick = 500ms real = 1 minuto de jogo. 1 dia de jogo = 1440 ticks.
// Esses valores fazem a fome chegar ao máximo em ~1 dia de jogo.
const (
	hungerDecayRate  = 1.0 / 1440  // 0 → 1.0 em 1 dia de jogo
	fatigueWorkRate  = 1.0 / 960   // 0 → 1.0 em ~16h de jogo (trabalhando)
	fatigueSleepRate = -(1.0 / 360) // 1.0 → 0 em ~6h de jogo (dormindo)
	fatigueIdleRate  = 1.0 / 1440  // sobe devagar quando ocioso
)

// tick atualiza as necessidades com base na atividade atual.
// É privado (letra minúscula) — só o pacote npc chama isso.
func (n *Needs) tick(activity Activity) {
	n.Hunger = clamp(n.Hunger+hungerDecayRate, 0, 1)

	switch activity {
	case ActivitySleep:
		n.Fatigue = clamp(n.Fatigue+fatigueSleepRate, 0, 1)
	case ActivityWork, ActivityPatrol:
		n.Fatigue = clamp(n.Fatigue+fatigueWorkRate, 0, 1)
	default:
		n.Fatigue = clamp(n.Fatigue+fatigueIdleRate, 0, 1)
	}
}

// NPC representa um personagem não-jogável vivendo no mundo.
//
// CurrentZoneID é a localização atual — uma chave do WorldMap.
// lastPeriod rastreia o período anterior para detectar quando o período muda.
// Não exportamos lastPeriod (minúsculo) porque é detalhe interno de estado.
type NPC struct {
	ID              string
	Name            string
	CurrentZoneID   string
	Schedule        Schedule
	Needs           Needs
	currentActivity Activity
	lastPeriod      string
}

// ProcessTick é chamado a cada tick do game loop.
// period é o período atual do jogo: "Manha", "Tarde", "Noite", "Madrugada".
//
// O NPC faz duas coisas por tick:
//  1. Atualiza as necessidades com base na atividade atual.
//  2. Se o período mudou, consulta a agenda e se move para a zona correta.
func (n *NPC) ProcessTick(period string) {
	// Atualiza necessidades com a atividade ATUAL (antes de trocar de período).
	n.Needs.tick(n.currentActivity)

	// Sem mudança de período → nada mais a fazer.
	if period == n.lastPeriod {
		return
	}
	n.lastPeriod = period

	// Busca o que fazer neste período.
	// Nota Go: o segundo valor do map lookup é um bool indicando se a chave existe.
	// Isso evita panic (panic seria: acessar chave inexistente em mapa nil).
	entry, ok := n.Schedule[period]
	if !ok {
		return
	}

	n.CurrentZoneID = entry.ZoneID
	n.currentActivity = entry.Activity
}

// NPCState é um snapshot read-only do estado atual de um NPC.
// Exportamos um struct separado (em vez de tornar os campos do NPC públicos)
// pelo mesmo motivo que World tem Snapshot: quem lê não precisa do mutex,
// e quem escreve não expõe acidentalmente campos internos de controle.
type NPCState struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	ZoneID   string   `json:"zone_id"`
	Activity Activity `json:"activity"`
	Needs    Needs    `json:"needs"`
}

// State retorna um snapshot do estado atual do NPC.
// Não precisa de mutex aqui porque NPC é sempre acessado dentro do
// mutex do World (em ProcessTick e em leituras via Snapshot).
func (n *NPC) State() NPCState {
	return NPCState{
		ID:       n.ID,
		Name:     n.Name,
		ZoneID:   n.CurrentZoneID,
		Activity: n.currentActivity,
		Needs:    n.Needs,
	}
}

// clamp garante que v fique entre min e max.
// Função auxiliar privada — usada apenas dentro do pacote npc.
// Go não tem função clamp na stdlib para float64 (só a partir do Go 1.21 tem em math).
func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
