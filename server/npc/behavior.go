package npc

import (
	"log/slog"

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
	if npc.CurrentActivity == activity && npc.CurrentZone == location {
		return
	}

	previous := npc.CurrentActivity
	npc.CurrentActivity = activity
	npc.CurrentZone = location
	npc.ActivityStartedAt = gameTime.Time

	slog.Info("npc transitioned",
		"npc", npc.Id,
		"from", previous,
		"to", activity,
		"location", location,
		"hunger", npc.Needs.Hunger,
		"fatigue", npc.Needs.Fatigue,
	)
}
