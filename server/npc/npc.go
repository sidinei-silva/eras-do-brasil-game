package npc

import "time"

type Role string

const (
	RoleMerchant   Role = "merchant"
	RoleQuestGiver Role = "quest_giver"
	RoleGuard      Role = "guard"
	RoleVillager   Role = "villager"
)

type Npc struct {
	Id                string           // Identificador único do NPC
	Name              string           // Nome do NPC
	Role              Role             // Função ou profissão do NPC (ex.: comerciante, entregador de missão, guarda)
	CurrentBlock      string           // Bloco atual onde o NPC está localizado
	CurrentActivity   Activity         // Atividade atual em que o NPC está envolvido
	Description       string           // Descrição da aparência e personalidade do NPC
	Backstory         string           // História de fundo do NPC
	Schedule          []ScheduleAction // Agenda diária das atividades do NPC
	Needs             Need             // Necessidades atuais do NPC
	HomePoi           string           // Local onde o NPC mora para dormir
	EatingPoi         string           // Local onde o NPC come
	ActivityStartedAt time.Time        // Hora em que a atividade atual começou
	NeedsWeight       NeedWeight       // Pesos para cada necessidade, usados na decisão de atividades
	CurrentPoi        string           // Ponto de interesse onde o npc está atualmente
	validPois         map[string]struct{}
}

func NewNpc(id string, name string, role Role, currentBlock string, description string, backstory string, schedule []ScheduleAction, homePoi string, eatingPoi string, validPois []string, needsWeight NeedWeight) *Npc {
	poiSet := make(map[string]struct{}, len(validPois))
	for _, poiId := range validPois {
		poiSet[poiId] = struct{}{}
	}

	return &Npc{
		Id:              id,
		Name:            name,
		Role:            role,
		CurrentBlock:    currentBlock,
		CurrentActivity: "idle", // A atividade padrão é idle até que a agenda seja avaliada
		Description:     description,
		Backstory:       backstory,
		Schedule:        schedule,
		Needs: Need{
			Hunger:     0,
			Fatigue:    0,
			Loneliness: 0,
		},
		HomePoi:     homePoi,
		EatingPoi:   eatingPoi,
		NeedsWeight: needsWeight,
		CurrentPoi:  "",
		validPois:   poiSet,
	}
}

func (npc *Npc) hasValidPoi(poiId string) bool {
	if poiId == "" {
		return true
	}

	_, exists := npc.validPois[poiId]
	return exists
}
