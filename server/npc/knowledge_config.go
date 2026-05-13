package npc

import (
	"encoding/json"
	"log/slog"
	"os"
)

type KnowledgeConfig struct {
	ExpirationHours           map[KnowledgeType]int `json:"expirationHours"`           // Horas para expiração de cada tipo de conhecimento
	GossipCooldownHours       int                   `json:"gossipCooldownHours"`       // Horas de cooldown para fofocas entre NPCs
	MaxKnowledgeBaseSize      int                   `json:"maxKnowledgeBaseSize"`      // Tamanho máximo da base de conhecimento de um NPC
	GossipMinExchangePerEvent int                   `json:"gossipMinExchangePerEvent"` // Número mínimo de conhecimentos trocados em um evento de fofoca
	GossipMaxExchangePerEvent int                   `json:"gossipMaxExchangePerEvent"` // Número máximo de conhecimentos trocados em um evento de fofoca
}

func LoadKnowledgeConfig() (*KnowledgeConfig, error) {
	// Implementação para carregar a configuração de conhecimento, possivelmente de um arquivo JSON

	filePath := os.Getenv("KNOWLEDGE_CONFIG_FILE")

	if filePath == "" {
		// Caminho padrão quando a variável KNOWLEDGE_CONFIG_FILE não estiver definida.
		filePath = "./data/knowledge_config.json"
	}

	jsonFile, err := os.ReadFile(filePath)

	if err != nil {
		slog.Error("Falha ao carregar arquivo", "err", err)
		return nil, err
	}

	var config KnowledgeConfig
	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		slog.Error("Falha ao unmarshal JSON", "err", err)
		return nil, err
	}

	// verificar se valores é 0 ou negativo causar um os.Exit, pois não faz sentido ter valores negativos ou 0 para essas configurações
	if config.GossipCooldownHours <= 0 {
		slog.Error("Valor inválido para gossipCooldownHours", "value", config.GossipCooldownHours)
		os.Exit(1)
	}
	if config.MaxKnowledgeBaseSize <= 0 {
		slog.Error("Valor inválido para maxKnowledgeBaseSize", "value", config.MaxKnowledgeBaseSize)
		os.Exit(1)
	}
	if config.GossipMinExchangePerEvent <= 0 {
		slog.Error("Valor inválido para gossipMinExchangePerEvent", "value", config.GossipMinExchangePerEvent)
		os.Exit(1)
	}
	if config.GossipMaxExchangePerEvent <= 0 {
		slog.Error("Valor inválido para gossipMaxExchangePerEvent", "value", config.GossipMaxExchangePerEvent)
		os.Exit(1)
	}

	// Verificar se os tipos de conhecimento possuem um valor de expiração definido
	for knowledgeType, hours := range config.ExpirationHours {
		if hours <= 0 {
			slog.Error("Valor inválido para expiração de conhecimento", "knowledgeType", knowledgeType, "value", hours)
			os.Exit(1)
		}
	}

	return &config, nil
}
