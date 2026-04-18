package npc

import (
	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/domain/world"
)

type Activity string
type Role string

// TODO: pesquisar um bom valor para a taxa de decaimento de fome e energia, considerando o avanço do tempo no jogo (ex: 1 ponto de fome a cada 10 minutos de jogo, 1 ponto de energia a cada 30 minutos de jogo, etc)
const (
	HugerDecayRate  = 1 / 1440.0 // 1 por minuto, considerando que o tempo do jogo avança 60 vezes mais rápido que o tempo real (1 minuto de jogo = 1 segundo real)
	EnergyDecayRate = 1 / 1440.0 // 1 por 12 minutos, considerando que o tempo do jogo avança 60 vezes mais rápido que o tempo real (12 minutos de jogo = 1 segundo real)
)

const (
	RoleMerchant   Role = "merchant"
	RoleQuestGiver Role = "quest_giver"
	RoleGuard      Role = "guard"
	RoleVillager   Role = "villager"
)
const (
	ActivityEating   Activity = "eating"
	ActivityWalking  Activity = "walking"
	ActivityIdle     Activity = "idle"
	ActivitySleeping Activity = "sleeping"
	ActivityWorking  Activity = "working"
)

type NpcScheduleAction struct {
	Activity Activity          // The activity the NPC will perform
	Location string            // The location where the activity will take place
	Period   world.PeriodOfDay // The period of the day when the activity will occur
}

type NpcNeed struct {
	Hunger float64 // Hunger level (0 to 100)
	Energy float64 // Energy level (0 to 100)
}

type Npc struct {
	ID              string              // Unique identifier for the NPC
	Name            string              // Name of the NPC
	Role            Role                // Role or profession of the NPC (e.g., merchant, quest giver, guard)
	CurrentZone     string              // Current zone where the NPC is located
	CurrentActivity Activity            // Current activity the NPC is engaged in
	Description     string              // Description of the NPC's appearance and personality
	Backstory       string              // Background story of the NPC
	Schedule        []NpcScheduleAction // Daily schedule of the NPC's activities
	Needs           NpcNeed             // Current needs of the NPC
}

func NewNpc(id string, name string, role Role, currentZone string, description string, backstory string, schedule []NpcScheduleAction) *Npc {
	return &Npc{
		ID:              id,
		Name:            name,
		Role:            role,
		CurrentZone:     currentZone,
		CurrentActivity: "idle", // Default activity is idle until the schedule is evaluated
		Description:     description,
		Backstory:       backstory,
		Schedule:        schedule,
		Needs: NpcNeed{
			Hunger: 0,
			Energy: 100,
		},
	}
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

// Get the NPC's current activity and location based on the game time, the NPC's schedule and need levels
// TODO: Colocar nos npcs os locais favoritos baseado nos dados do npc e usar isso para escolher a localização quando a atividade for "walking" ou "idle"
func (npc *Npc) CurrentActivityAndLocation(gameTime world.GameTime) {
	const (
		hungerCritical = 80.0
		energyLow      = 20.0
	)

	// Need-based overrides have priority over schedule.
	if npc.Needs.Hunger >= hungerCritical {
		action := npc.getNpcScheduleActionEating()
		npc.CurrentActivity = action.Activity
		npc.CurrentZone = action.Location
	}

	if npc.Needs.Energy <= energyLow {
		action := npc.getNpcScheduleActionSleeping()
		npc.CurrentActivity = action.Activity
		npc.CurrentZone = action.Location
	}

	for _, action := range npc.Schedule {
		if action.Period == gameTime.PeriodOfDay {
			location := action.Location
			if location == "" {
				location = npc.CurrentZone
			}

			npc.CurrentActivity = action.Activity
			npc.CurrentZone = location
			return
		}
	}

	npc.CurrentActivity = ActivityIdle
}

func (npc *Npc) getNpcScheduleActionEating() NpcScheduleAction {
	for _, action := range npc.Schedule {
		if action.Activity == ActivityEating {
			return action
		}
	}

	return NpcScheduleAction{}
}

func (npc *Npc) getNpcScheduleActionSleeping() NpcScheduleAction {
	for _, action := range npc.Schedule {
		if action.Activity == ActivitySleeping {
			return action
		}
	}

	return NpcScheduleAction{}
}
