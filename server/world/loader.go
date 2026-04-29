package world

import (
	"encoding/json"
	"log/slog"
	"os"
	"slices"
)

type TemplateConnection struct {
	Terrain       string `json:"terrain"`
	ToBlockId     string `json:"toBlockId"`
	TravelMinutes int    `json:"travelMinutes"`
}

type TemplateLevelRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type TemplateBlock struct {
	Id          string               `json:"id"`
	Name        string               `json:"name"`
	Type        string               `json:"type" validate:"oneof=urban dense_forest mountain water ruins wasteland"`
	Description string               `json:"description"`
	Connections []TemplateConnection `json:"connections"`
	LevelRange  TemplateLevelRange   `json:"levelRange"`
	Region      Region               `json:"region"`
	Tags        []string             `json:"tags"`
}

type Data struct {
	Blocks []TemplateBlock `json:"blocks"`
}

func LoadBlocksFromFile() ([]*Block, error) {
	filePath := os.Getenv("BLOCKS_FILE")

	if filePath == "" {
		// Caminho padrão quando a variável BLOCKS_FILE não estiver definida.
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

	blocks := make([]*Block, 0, len(data.Blocks))

	for _, blockData := range data.Blocks {
		block := NewBlock(
			blockData.Id,
			blockData.Name,
			BlockType(blockData.Type),
			blockData.Description,
			LevelRange{
				Min: blockData.LevelRange.Min,
				Max: blockData.LevelRange.Max,
			},
			blockData.Region,
			blockData.Tags,
		)

		for _, connData := range blockData.Connections {

			checkIfExist := slices.ContainsFunc(data.Blocks, func(item TemplateBlock) bool {
				return item.Id == connData.ToBlockId
			})

			if !checkIfExist {
				slog.Error("Conexão para bloco inexistente", "fromBlockId", blockData.Id, "toBlockId", connData.ToBlockId)
				os.Exit(1)
				continue
			}

			connection := Connection{
				Terrain:       connData.Terrain,
				ToBlockId:     connData.ToBlockId,
				TravelMinutes: connData.TravelMinutes,
			}
			block.Connections = append(block.Connections, connection)
		}

		blocks = append(blocks, block)
	}

	return blocks, nil
}
