package npc

type Role string

const (
	RoleMerchant   Role = "merchant"
	RoleQuestGiver Role = "quest_giver"
	RoleGuard      Role = "guard"
	RoleVillager   Role = "villager"
)

type Npc struct {
	Id              string           // Unique identifier for the NPC
	Name            string           // Name of the NPC
	Role            Role             // Role or profession of the NPC (e.g., merchant, quest giver, guard)
	CurrentZone     string           // Current zone where the NPC is located
	CurrentActivity Activity         // Current activity the NPC is engaged in
	Description     string           // Description of the NPC's appearance and personality
	Backstory       string           // Background story of the NPC
	Schedule        []ScheduleAction // Daily schedule of the NPC's activities
	Needs           Need             // Current needs of the NPC
}

func NewNpc(id string, name string, role Role, currentZone string, description string, backstory string, schedule []ScheduleAction) *Npc {
	return &Npc{
		Id:              id,
		Name:            name,
		Role:            role,
		CurrentZone:     currentZone,
		CurrentActivity: "idle", // Default activity is idle until the schedule is evaluated
		Description:     description,
		Backstory:       backstory,
		Schedule:        schedule,
		Needs: Need{
			Hunger: 0,
			Energy: 100,
		},
	}
}
