package npc

import "github.com/sidinei-silva/eras-do-brasil-game/server/world"

// Get the NPC's current activity and location based on the game time, the NPC's schedule and need levels
// TODO: Colocar nos npcs os locais favoritos baseado nos dados do npc e usar isso para escolher a localização quando a atividade for "walking" ou "idle"
func (npc *Npc) CurrentActivityAndLocation(gameTime world.GameTime) {
	const (
		hungerCritical = 80.0
		energyLow      = 20.0
	)

	// Need-based overrides have priority over schedule.
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

func (npc *Npc) getScheduleActionEating() ScheduleAction {
	for _, action := range npc.Schedule {
		if action.Activity == ActivityEating {
			return action
		}
	}

	return ScheduleAction{
		Activity: ActivityIdle,
		Location: npc.CurrentZone,
	}
}

func (npc *Npc) getScheduleActionSleeping() ScheduleAction {
	for _, action := range npc.Schedule {
		if action.Activity == ActivitySleeping {
			return action
		}
	}

	return ScheduleAction{
		Activity: ActivityIdle,
		Location: npc.CurrentZone,
	}
}
