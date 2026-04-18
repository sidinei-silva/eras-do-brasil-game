package data

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/domain/npc"
	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/domain/world"
)

type NpcScheduleActionDTO struct {
	Activity string `json:"activity"`
	Location string `json:"location"`
	Period   string `json:"period"`
}

type NpcTemplateDTO struct {
	Backstory   string                 `json:"backstory"`
	CurrentZone string                 `json:"currentZone"`
	Description string                 `json:"description"`
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Role        string                 `json:"role"`
	Schedule    []NpcScheduleActionDTO `json:"schedule"`
}

type NpcData struct {
	Npcs []NpcTemplateDTO `json:"npcs"`
}

func LoadNpcsFromFile() (map[string]*npc.Npc, error) {
	filePath := os.Getenv("NPCS_FILE")

	if filePath == "" {
		// fallback para quando o servidor é iniciado em server/cmd/game
		filePath = "../../data/npcs.json"
	}

	jsonFile, err := os.ReadFile(filePath)

	if err != nil {
		slog.Error("Falha ao carregar arquivo", "err", err)
		return nil, err
	}

	var data NpcData
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		slog.Error("Falha ao unmarshal JSON", "err", err)
		return nil, err
	}

	slog.Info("NPCs carregados com sucesso", "count", len(data.Npcs))

	npcs := make(map[string]*npc.Npc)

	for _, npcData := range data.Npcs {
		npc := npc.NewNpc(
			npcData.Id,
			npcData.Name,
			npc.Role(npcData.Role),
			npcData.CurrentZone,
			npcData.Description,
			npcData.Backstory,
			func() []npc.NpcScheduleAction {
				schedule := make([]npc.NpcScheduleAction, len(npcData.Schedule))
				for i, action := range npcData.Schedule {
					schedule[i] = npc.NpcScheduleAction{
						Activity: npc.Activity(action.Activity),
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
