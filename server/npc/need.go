package npc

import "github.com/sidinei-silva/eras-do-brasil-game/server/world"

// TODO: pesquisar um bom valor para a taxa de decaimento de fome e energia, considerando o avanço do tempo no jogo (ex: 1 ponto de fome a cada 10 minutos de jogo, 1 ponto de energia a cada 30 minutos de jogo, etc)
const (
	HugerDecayRate   = 1 / 1440.0
	FatigueDecayRate = 1 / 1440.0
)

type Need struct {
	Hunger     float64 // Nível de fome (0 a 100)
	Fatigue    float64 // Nível de fadiga (0 a 100)
	Loneliness float64 // Nível de solidão (0 a 100)
}

// TODO: Melhorar o balanceamento de multiplicador de necessidade baseado na atividade e no tempo do jogo (ex: se for noite, dormir aumenta mais fatigue, comer aumenta mais a fome, etc)
func (npc *Npc) UpdateNeeds(gameTime world.GameTime) {

	switch npc.CurrentActivity {
	case ActivityWorking:
		npc.Needs.Hunger += HugerDecayRate
		npc.Needs.Fatigue += FatigueDecayRate
	case ActivitySleeping:
		npc.Needs.Hunger += HugerDecayRate / 2
		npc.Needs.Fatigue -= FatigueDecayRate * 2
	case ActivityIdle:
		npc.Needs.Hunger += HugerDecayRate / 2
		npc.Needs.Fatigue += FatigueDecayRate / 2
	case ActivityEating:
		npc.Needs.Hunger -= HugerDecayRate * 2
		npc.Needs.Fatigue += FatigueDecayRate

	}

}
