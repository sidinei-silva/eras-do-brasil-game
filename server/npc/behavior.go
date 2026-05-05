package npc

import (
	"github.com/sidinei-silva/eras-do-brasil-game/server/world"
)

// IsDiscreteActivity responde se a atividade atual tem duração fixa
// definida (almoço, sono, refeição). Atividades sem duração são
// contínuas (working, idle, walking).
func (npc *Npc) IsDiscreteActivity() bool {
	_, exists := activityDuration[npc.CurrentActivity]
	return exists
}

// IsActivityComplete responde se já passou tempo suficiente desde que
// a atividade discreta atual começou.
func (npc *Npc) IsActivityComplete(gameTime world.GameTime) bool {
	duration := getActivityDuration(npc.CurrentActivity)
	elapsed := gameTime.Time.Sub(npc.ActivityStartedAt)
	return elapsed >= duration
}

// TransitionTo muda a atividade e zona do NPC, registrando o momento
// da troca. É chamado pelo Manager quando a decisão muda a atividade.
func (npc *Npc) TransitionTo(activity Activity, location string, gameTime world.GameTime) {
	if npc.CurrentActivity == activity && npc.CurrentBlock == location {
		return
	}

	npc.CurrentActivity = activity
	npc.CurrentBlock = location
	npc.ActivityStartedAt = gameTime.Time

}
