package npc

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed data/mata-costeira-npcs.json
var mataCosteiraNPCsSrc []byte

type scheduleEntryJSON struct {
	ZoneID   string `json:"zone_id"`
	Activity string `json:"activity"`
}

type npcJSON struct {
	ID        string                       `json:"id"`
	Name      string                       `json:"name"`
	SpawnZone string                       `json:"spawn_zone"`
	Schedule  map[string]scheduleEntryJSON `json:"schedule"`
}

type npcDataJSON struct {
	NPCs []npcJSON `json:"npcs"`
}

func LoadMataCosteira() (*Registry, error) {
	var data npcDataJSON
	if err := json.Unmarshal(mataCosteiraNPCsSrc, &data); err != nil {
		return nil, fmt.Errorf("loading mata-costeira-npcs.json: %w", err)
	}

	r := NewRegistry()

	for _, nj := range data.NPCs {
		schedule := make(Schedule, len(nj.Schedule))
		for period, ej := range nj.Schedule {
			schedule[period] = ScheduleEntry{
				ZoneID:   ej.ZoneID,
				Activity: Activity(ej.Activity),
			}
		}

		r.npcs[nj.ID] = &NPC{
			ID:              nj.ID,
			Name:            nj.Name,
			CurrentZoneID:   nj.SpawnZone,
			Schedule:        schedule,
			currentActivity: ActivityIdle,
		}
	}

	return r, nil
}
