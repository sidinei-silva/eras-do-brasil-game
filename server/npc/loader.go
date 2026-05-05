package npc

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/sidinei-silva/eras-do-brasil-game/server/config"
)

type NeedsWeightDTO struct {
	Hunger     float64 `json:"hunger"`
	Fatigue    float64 `json:"fatigue"`
	Loneliness float64 `json:"loneliness"`
	Schedule   float64 `json:"schedule"`
}

type ScheduleActionDTO struct {
	Activity  string `json:"activity"`
	Location  string `json:"location"`
	StartHour int    `json:"startHour"`
	EndHour   int    `json:"endHour"`
}

type TemplateDTO struct {
	Backstory      string              `json:"backstory"`
	CurrentZone    string              `json:"currentZone"`
	Description    string              `json:"description"`
	Id             string              `json:"id"`
	Name           string              `json:"name"`
	Role           string              `json:"role"`
	Schedule       []ScheduleActionDTO `json:"schedule"`
	HomeLocation   string              `json:"homeLocation"`
	EatingLocation string              `json:"eatingLocation"`
	NeedsWeight    NeedsWeightDTO      `json:"needsWeight"`
}

type Data struct {
	Npcs []TemplateDTO `json:"npcs"`
}

func LoadNpcsFromFile() ([]*Npc, error) {
	filePath := os.Getenv("NPCS_FILE")

	if filePath == "" {
		// Caminho padrão quando a variável NPCS_FILE não estiver definida.
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

	npcs := make([]*Npc, 0, len(data.Npcs))

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
						Activity:  Activity(action.Activity),
						Location:  action.Location,
						StartHour: action.StartHour,
						EndHour:   action.EndHour,
					}
				}
				return schedule
			}(),
			npcData.HomeLocation,
			npcData.EatingLocation,
			NeedWeight{
				Hunger:     npcData.NeedsWeight.Hunger,
				Fatigue:    npcData.NeedsWeight.Fatigue,
				Loneliness: npcData.NeedsWeight.Loneliness,
				Schedule:   npcData.NeedsWeight.Schedule,
			},
		)
		if config.Log.WorldLoading {
			slog.Debug("npc carregado", "id", npc.Id, "name", npc.Name, "role", npc.Role, "zone", npc.CurrentZone)
		}
		npcs = append(npcs, npc)
	}

	return npcs, nil
}
