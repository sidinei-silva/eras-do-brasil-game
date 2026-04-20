package npc

import "github.com/sidinei-silva/eras-do-brasil-game/server/world"

type Activity string

const (
	ActivityEating   Activity = "eating"
	ActivityWalking  Activity = "walking"
	ActivityIdle     Activity = "idle"
	ActivitySleeping Activity = "sleeping"
	ActivityWorking  Activity = "working"
)

type ScheduleAction struct {
	Activity Activity          // The activity the NPC will perform
	Location string            // The location where the activity will take place
	Period   world.PeriodOfDay // The period of the day when the activity will occur
}
