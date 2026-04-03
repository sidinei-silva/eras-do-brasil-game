package npc

type Activity string

const (
	ActivitySleep  Activity = "sleep"
	ActivityWork   Activity = "work"
	ActivityEat    Activity = "eat"
	ActivityPatrol Activity = "patrol"
	ActivityIdle   Activity = "idle"
)

type Needs struct {
	Hunger  float64
	Fatigue float64
}

const (
	hungerDecayRate  = 1.0 / 1440
	fatigueWorkRate  = 1.0 / 960
	fatigueSleepRate = -(1.0 / 360)
	fatigueIdleRate  = 1.0 / 1440
)

func (n *Needs) tick(activity Activity) {
	n.Hunger = clamp(n.Hunger+hungerDecayRate, 0, 1)

	switch activity {
	case ActivitySleep:
		n.Fatigue = clamp(n.Fatigue+fatigueSleepRate, 0, 1)
	case ActivityWork, ActivityPatrol:
		n.Fatigue = clamp(n.Fatigue+fatigueWorkRate, 0, 1)
	default:
		n.Fatigue = clamp(n.Fatigue+fatigueIdleRate, 0, 1)
	}
}

type NPC struct {
	ID              string
	Name            string
	CurrentZoneID   string
	Schedule        Schedule
	Needs           Needs
	currentActivity Activity
	lastPeriod      string
}

func (n *NPC) ProcessTick(period string) {
	n.Needs.tick(n.currentActivity)

	if period == n.lastPeriod {
		return
	}
	n.lastPeriod = period

	entry, ok := n.Schedule[period]
	if !ok {
		return
	}

	n.CurrentZoneID = entry.ZoneID
	n.currentActivity = entry.Activity
}

type NPCState struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	ZoneID   string   `json:"zone_id"`
	Activity Activity `json:"activity"`
	Needs    Needs    `json:"needs"`
}

func (n *NPC) State() NPCState {
	return NPCState{
		ID:       n.ID,
		Name:     n.Name,
		ZoneID:   n.CurrentZoneID,
		Activity: n.currentActivity,
		Needs:    n.Needs,
	}
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
