package npc

import "github.com/sidinei-silva/eras-do-brasil-game/server/world"

// Obtém a atividade e a localização atuais do NPC com base no horário do jogo, na agenda do NPC e nos níveis de necessidade
// TODO: Colocar nos npcs os locais favoritos baseado nos dados do npc e usar isso para escolher a localização quando a atividade for "walking" ou "idle"
func (npc *Npc) CurrentActivityAndLocation(gameTime world.GameTime) {

	// Isso será trocado por uma logica de scoring das necessidades do npc baseado em peso
	/*const (
		hungerCritical = 80.0
		energyLow      = 20.0
	)

	 if npc.Needs.Hunger >= hungerCritical {
		action := npc.getScheduleActionEating()
		npc.CurrentActivity = action.Activity
		npc.CurrentZone = action.Location
		return
	}

	if npc.Needs.Energy <= energyLow {
		action := npc.getScheduleActionSleeping()
		npc.CurrentActivity = action.Activity
		npc.CurrentZone = action.Location
		return
	} */

	actionActive, found := npc.ActiveScheduleAt(gameTime.Time.Hour())

	if found {
		npc.CurrentActivity = actionActive.Activity
		npc.CurrentZone = actionActive.Location
		return
	}

	npc.CurrentActivity = ActivityIdle
}
