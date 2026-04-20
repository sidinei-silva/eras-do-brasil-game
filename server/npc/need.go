package npc

import "github.com/sidinei-silva/eras-do-brasil-game/server/world"

// TODO: pesquisar um bom valor para a taxa de decaimento de fome e energia, considerando o avanço do tempo no jogo (ex: 1 ponto de fome a cada 10 minutos de jogo, 1 ponto de energia a cada 30 minutos de jogo, etc)
const (
	HugerDecayRate  = 1 / 1440.0 // 1 por minuto, considerando que o tempo do jogo avança 60 vezes mais rápido que o tempo real (1 minuto de jogo = 1 segundo real)
	EnergyDecayRate = 1 / 1440.0 // 1 por 12 minutos, considerando que o tempo do jogo avança 60 vezes mais rápido que o tempo real (12 minutos de jogo = 1 segundo real)
)

type Need struct {
	Hunger float64 // Hunger level (0 to 100)
	Energy float64 // Energy level (0 to 100)
}

// Update the NPC's hunger and energy levels in tick
// If working hunger increases and energy decreases.
// If sleeping hunger decreases and energy increases.
// If idle hunger increases slightly and energy increases slightly.
// If eating hunger decreases and energy increases.
// TODO: Melhorar o balanceamento de multiplicador de necessidade baseado na atividade e no tempo do jogo (ex: se for noite, dormir aumenta mais energia, comer aumenta mais a fome, etc)
func (npc *Npc) UpdateNeeds(gameTime world.GameTime) {

	switch npc.CurrentActivity {
	case ActivityWorking:
		npc.Needs.Hunger += HugerDecayRate
		npc.Needs.Energy -= EnergyDecayRate
	case ActivitySleeping:
		npc.Needs.Hunger -= HugerDecayRate / 2
		npc.Needs.Energy += EnergyDecayRate * 2
	case ActivityIdle:
		npc.Needs.Hunger += HugerDecayRate / 2
		npc.Needs.Energy += EnergyDecayRate / 2
	case ActivityEating:
		npc.Needs.Hunger -= HugerDecayRate * 2
		npc.Needs.Energy += EnergyDecayRate

	}

	// Ensure hunger and energy levels are within bounds (0 to 100)
	if npc.Needs.Hunger < 0 {
		npc.Needs.Hunger = 0
	} else if npc.Needs.Hunger > 100 {
		npc.Needs.Hunger = 100
	}

	if npc.Needs.Energy < 0 {
		npc.Needs.Energy = 0
	} else if npc.Needs.Energy > 100 {
		npc.Needs.Energy = 100
	}
}
