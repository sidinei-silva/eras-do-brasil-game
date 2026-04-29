package npc

type Activity string

type ScheduleAction struct {
	Activity  Activity // A atividade que o NPC irá realizar
	Location  string   // O local onde a atividade acontecerá
	StartHour int      // Hora de início da atividade (0-23)
	EndHour   int      // Hora de término da atividade (0-23)
}

func (npc *Npc) ActiveScheduleAt(hour int) (ScheduleAction, bool) {
	for _, action := range npc.Schedule {
		if hour >= action.StartHour && hour < action.EndHour {
			return action, true
		}
	}
	return ScheduleAction{}, false
}
