package npc

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/sidinei-silva/eras-do-brasil-game/server/world"
)

type ScheduleActionDTO struct {
	Activity string `json:"activity"`
	Location string `json:"location"`
	Period   string `json:"period"`
}

type TemplateDTO struct {
	Backstory   string              `json:"backstory"`
	CurrentZone string              `json:"currentZone"`
	Description string              `json:"description"`
	Id          string              `json:"id"`
	Name        string              `json:"name"`
	Role        string              `json:"role"`
	Schedule    []ScheduleActionDTO `json:"schedule"`
}

type Data struct {
	Npcs []TemplateDTO `json:"npcs"`
}

func LoadNpcsFromFile() (map[string]*Npc, error) {
	filePath := os.Getenv("NPCS_FILE")

	if filePath == "" {
		// fallback para quando o servidor é iniciado em server/cmd/game
		filePath = "./data/npcs.json"
	}

	jsonFile, err := os.ReadFile(filePath)

	if err != nil {
		slog.Error("Falha ao carregar arquivo", "err", err)
		return nil, err
	}

	var data Data
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		slog.Error("Falha ao unmarshal JSON", "err", err)
		return nil, err
	}

	npcs := make(map[string]*Npc)

	for _, npcData := range data.Npcs {
		npc := NewNpc(
			npcData.Id,
			npcData.Name,
			Role(npcData.Role),
			npcData.CurrentZone,
			npcData.Description,
			npcData.Backstory,
			func() []ScheduleAction {
				schedule := make([]ScheduleAction, len(npcData.Schedule))
				for i, action := range npcData.Schedule {
					schedule[i] = ScheduleAction{
						Activity: Activity(action.Activity),
						Location: action.Location,
						Period:   world.PeriodOfDay(action.Period),
					}
				}
				return schedule
			}(),
		)
		npcs[npc.Name] = npc
	}

	return npcs, nil
}
