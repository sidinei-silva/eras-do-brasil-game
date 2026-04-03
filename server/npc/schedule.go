package npc

type ScheduleEntry struct {
	ZoneID   string
	Activity Activity
}

type Schedule map[string]ScheduleEntry
