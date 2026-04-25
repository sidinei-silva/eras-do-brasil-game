package npc

type Role string

const (
	RoleMerchant   Role = "merchant"
	RoleQuestGiver Role = "quest_giver"
	RoleGuard      Role = "guard"
	RoleVillager   Role = "villager"
)

type Npc struct {
	Id              string           // Identificador único do NPC
	Name            string           // Nome do NPC
	Role            Role             // Função ou profissão do NPC (ex.: comerciante, entregador de missão, guarda)
	CurrentZone     string           // Zona atual onde o NPC está localizado
	CurrentActivity Activity         // Atividade atual em que o NPC está envolvido
	Description     string           // Descrição da aparência e personalidade do NPC
	Backstory       string           // História de fundo do NPC
	Schedule        []ScheduleAction // Agenda diária das atividades do NPC
	Needs           Need             // Necessidades atuais do NPC
	HomeLocation    string           // Local onde o NPC mora para dormir
	EatingLocation  string           // Local onde o NPC come
}

func NewNpc(id string, name string, role Role, currentZone string, description string, backstory string, schedule []ScheduleAction, homeLocation string, eatingLocation string) *Npc {
	return &Npc{
		Id:              id,
		Name:            name,
		Role:            role,
		CurrentZone:     currentZone,
		CurrentActivity: "idle", // A atividade padrão é idle até que a agenda seja avaliada
		Description:     description,
		Backstory:       backstory,
		Schedule:        schedule,
		Needs: Need{
			Hunger: 0,
			Energy: 100,
		},
		HomeLocation:   homeLocation,
		EatingLocation: eatingLocation,
	}
}
