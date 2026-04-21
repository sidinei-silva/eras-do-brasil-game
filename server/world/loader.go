package world

import (
	"encoding/json"
	"log/slog"
	"os"
)

type TemplateConnection struct {
	Terrain           string `json:"terrain"`
	ToBlockId         string `json:"toBlockId"`
	TravelTimeMinutes int    `json:"travelTimeMinutes"`
}

type TemplateBlock struct {
	Id          string               `json:"id"`
	Name        string               `json:"name"`
	Type        string               `json:"type" validate:"oneof=urban dense_forest mountain water ruins wasteland"`
	Description string               `json:"description"`
	Connections []TemplateConnection `json:"connections"`
	LevelRange  [2]int               `json:"levelRange"`
	Region      Region               `json:"region"`
	Tags        []string             `json:"tags"`
}

type Data struct {
	Blocks []TemplateBlock `json:"blocks"`
}

func LoadBlocksFromFile() (map[string]*Block, error) {
	filePath := os.Getenv("BLOCKS_FILE")

	if filePath == "" {
		// fallback para quando o servidor é iniciado em server/cmd/game
		filePath = "./data/blocks.json"
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

	blocks := make(map[string]*Block)

	for _, blockData := range data.Blocks {
		block := NewBlock(
			blockData.Id,
			blockData.Name,
			blockData.Type,
			blockData.Description,
			blockData.LevelRange,
			blockData.Region,
			blockData.Tags,
		)

		for _, connData := range blockData.Connections {
			connection := Connection{
				Terrain:           connData.Terrain,
				ToBlockId:         connData.ToBlockId,
				TravelTimeMinutes: connData.TravelTimeMinutes,
			}
			block.Connections = append(block.Connections, connection)
		}

		blocks[block.Id] = block
	}

	return blocks, nil
}
