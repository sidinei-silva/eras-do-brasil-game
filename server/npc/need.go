package npc

// Taxa por hora de jogo
const (
	BaseHungerRate     = 12.5 // 100% em 8h de Idle
	BaseFatigueRate    = 4.0  // 100% em ~25h de Idle
	BaseLonelinessRate = 4.0  // 100% em 25h de Idle
)

type Need struct {
	Hunger     float64 // Nível de fome (0 a 100)
	Fatigue    float64 // Nível de fadiga (0 a 100)
	Loneliness float64 // Nível de solidão (0 a 100)
}

func (npc *Npc) ApplyDecay(tickHours float64) {
	intensity := getActivityIntensity(npc.CurrentActivity)

	npc.Needs.Hunger += BaseHungerRate * intensity * tickHours
	npc.Needs.Fatigue += BaseFatigueRate * intensity * tickHours
	npc.Needs.Loneliness += BaseLonelinessRate * tickHours
}

func (npc *Npc) ApplyActivityEffects() {
	effect := getActivityEffects(npc.CurrentActivity)
	npc.Needs.Hunger += effect.HungerDelta
	npc.Needs.Fatigue += effect.FatigueDelta
	npc.Needs.Loneliness += effect.LonelinessDelta

}

func (npc *Npc) ClampNeeds() {
	if npc.Needs.Hunger < 0 {
		npc.Needs.Hunger = 0
	} else if npc.Needs.Hunger > 100 {
		npc.Needs.Hunger = 100
	}

	if npc.Needs.Fatigue < 0 {
		npc.Needs.Fatigue = 0
	} else if npc.Needs.Fatigue > 100 {
		npc.Needs.Fatigue = 100
	}

	if npc.Needs.Loneliness < 0 {
		npc.Needs.Loneliness = 0
	} else if npc.Needs.Loneliness > 100 {
		npc.Needs.Loneliness = 100
	}
}
