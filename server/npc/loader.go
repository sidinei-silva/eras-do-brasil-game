package npc

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/sidinei-silva/eras-do-brasil-game/server/config"
	"github.com/sidinei-silva/eras-do-brasil-game/server/world"
)

type NeedsWeightDTO struct {
	Hunger     float64 `json:"hunger"`
	Fatigue    float64 `json:"fatigue"`
	Loneliness float64 `json:"loneliness"`
	Schedule   float64 `json:"schedule"`
}

type ScheduleActionDTO struct {
	Activity  string `json:"activity"`
	PoiId     string `json:"poiId"`
	StartHour int    `json:"startHour"`
	EndHour   int    `json:"endHour"`
}

type TemplateDTO struct {
	Backstory    string              `json:"backstory"`
	CurrentBlock string              `json:"currentBlock"`
	Description  string              `json:"description"`
	Id           string              `json:"id"`
	Name         string              `json:"name"`
	Role         string              `json:"role"`
	Schedule     []ScheduleActionDTO `json:"schedule"`
	HomePoi      string              `json:"homePoi"`
	EatingPoi    string              `json:"eatingPoi"`
	NeedsWeight  NeedsWeightDTO      `json:"needsWeight"`
}

type Data struct {
	Npcs []TemplateDTO `json:"npcs"`
}

func LoadNpcsFromFile() ([]*Npc, error) {
	blocks, err := world.LoadBlocksFromFile()
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

		checkLocationsInNpcExists(&npcData, blocks)

		npc := NewNpc(
			npcData.Id,
			npcData.Name,
			Role(npcData.Role),
			npcData.CurrentBlock,
			npcData.Description,
			npcData.Backstory,
			func() []ScheduleAction {
				schedule := make([]ScheduleAction, len(npcData.Schedule))
				for i, action := range npcData.Schedule {
					schedule[i] = ScheduleAction{
						Activity:  Activity(action.Activity),
						PoiId:     action.PoiId,
						StartHour: action.StartHour,
						EndHour:   action.EndHour,
					}
				}
				return schedule
			}(),
			npcData.HomePoi,
			npcData.EatingPoi,
			NeedWeight{
				Hunger:     npcData.NeedsWeight.Hunger,
				Fatigue:    npcData.NeedsWeight.Fatigue,
				Loneliness: npcData.NeedsWeight.Loneliness,
				Schedule:   npcData.NeedsWeight.Schedule,
			},
		)
		if config.Log.WorldLoading {
			slog.Debug("npc carregado", "id", npc.Id, "name", npc.Name, "role", npc.Role, "zone", npc.CurrentBlock)
		}
		npcs = append(npcs, npc)
	}

	return npcs, nil
}

func checkLocationsInNpcExists(npc *TemplateDTO, blocks []*world.Block) {

	// Verificar se CurrentBlock existe
	allBlocksIds := make(map[string]*world.Block)

	for _, block := range blocks {
		allBlocksIds[block.Id] = block
	}

	block, currentBlockExists := allBlocksIds[npc.CurrentBlock]
	if !currentBlockExists {
		slog.Error("Current block não encontrado para NPC", "npcId", npc.Id, "currentBlock", npc.CurrentBlock)
		panic("Current block não encontrado para NPC")
	}

	allPois := make(map[string]bool)
	for _, poi := range block.Pois {
		allPois[poi] = true
	}

	// Verificar se HomePoi existe
	if _, exists := allPois[npc.HomePoi]; !exists {
		slog.Error("Home POI não encontrado para NPC", "npcId", npc.Id, "homePoi", npc.HomePoi)
		panic("Home POI não encontrado para NPC")
	}

	// Verificar se EatingPoi existe
	if _, exists := allPois[npc.EatingPoi]; !exists {
		slog.Error("Eating POI não encontrado para NPC", "npcId", npc.Id, "eatingPoi", npc.EatingPoi)
		panic("Eating POI não encontrado para NPC")
	}

	// Verificar se os pois da schedule existem
	for _, action := range npc.Schedule {
		if _, exists := allPois[action.PoiId]; !exists {
			slog.Error("Schedule POI não encontrado para NPC", "npcId", npc.Id, "schedulePoi", action.PoiId)
			panic("Schedule POI não encontrado para NPC")
		}
	}

}
