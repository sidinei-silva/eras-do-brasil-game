package npc

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

type KnowledgeConfigDTO struct {
	ExpirationHours           map[KnowledgeType]int `json:"expirationHours"`           // Horas para expiração de cada tipo de conhecimento
	GossipCooldownHours       int                   `json:"gossipCooldownHours"`       // Horas de cooldown para fofocas entre NPCs
	MaxKnowledgeBaseSize      int                   `json:"maxKnowledgeBaseSize"`      // Tamanho máximo da base de conhecimento de um NPC
	GossipMinExchangePerEvent int                   `json:"gossipMinExchangePerEvent"` // Número mínimo de conhecimentos trocados em um evento de fofoca
	GossipMaxExchangePerEvent int                   `json:"gossipMaxExchangePerEvent"` // Número máximo de conhecimentos trocados em um evento de fofoca
}

type KnowledgeConfig struct {
	ExpirationHours           map[KnowledgeType]int // Horas para expiração de cada tipo de conhecimento
	GossipCooldownHours       int                   // Horas de cooldown para fofocas entre NPCs
	MaxKnowledgeBaseSize      int                   // Tamanho máximo da base de conhecimento de um NPC
	GossipMinExchangePerEvent int                   // Número mínimo de conhecimentos trocados em um evento de fofoca
	GossipMaxExchangePerEvent int                   // Número máximo de conhecimentos trocados em um evento de fofoca
}

func LoadKnowledgeConfig() (KnowledgeConfig, error) {
	// Implementação para carregar a configuração de conhecimento, possivelmente de um arquivo JSON

	filePath := os.Getenv("KNOWLEDGE_CONFIG_FILE")

	if filePath == "" {
		// Caminho padrão quando a variável KNOWLEDGE_CONFIG_FILE não estiver definida.
		filePath = "./data/knowledge_config.json"
	}

	jsonFile, err := os.ReadFile(filePath)

	if err != nil {
		slog.Error("Falha ao carregar arquivo", "err", err)
		return KnowledgeConfig{}, err
	}

	var config KnowledgeConfigDTO
	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		slog.Error("Falha ao unmarshal JSON", "err", err)
		return KnowledgeConfig{}, err
	}

	// verificar se valores é 0 ou negativo causar um os.Exit, pois não faz sentido ter valores negativos ou 0 para essas configurações
	if config.GossipCooldownHours <= 0 {
		return KnowledgeConfig{}, fmt.Errorf("gossipCooldownHours inválido: %d", config.GossipCooldownHours)
	}
	if config.MaxKnowledgeBaseSize <= 0 {
		return KnowledgeConfig{}, fmt.Errorf("maxKnowledgeBaseSize inválido: %d", config.MaxKnowledgeBaseSize)
	}
	if config.GossipMinExchangePerEvent <= 0 {
		return KnowledgeConfig{}, fmt.Errorf("gossipMinExchangePerEvent inválido: %d", config.GossipMinExchangePerEvent)
	}
	if config.GossipMaxExchangePerEvent <= 0 {
		return KnowledgeConfig{}, fmt.Errorf("gossipMaxExchangePerEvent inválido: %d", config.GossipMaxExchangePerEvent)
	}
	if config.GossipMaxExchangePerEvent < config.GossipMinExchangePerEvent {
		return KnowledgeConfig{}, fmt.Errorf("gossipMaxExchangePerEvent (%d) menor que min (%d)", config.GossipMaxExchangePerEvent, config.GossipMinExchangePerEvent)
	}

	// Verificar se os tipos de conhecimento possuem um valor de expiração definido
	for knowledgeType, hours := range config.ExpirationHours {
		if hours <= 0 {
			slog.Error("Valor inválido para expiração de conhecimento", "knowledgeType", knowledgeType, "value", hours)
			return KnowledgeConfig{}, err
		}
	}

	validTypes := map[KnowledgeType]bool{KnowledgeAvistamentoNPC: true}
	for kt := range config.ExpirationHours {
		if !validTypes[kt] {
			return KnowledgeConfig{}, fmt.Errorf("tipo de knowledge inválido em expirationHours: %s", kt)
		}
	}

	return KnowledgeConfig{
		ExpirationHours:           config.ExpirationHours,
		GossipCooldownHours:       config.GossipCooldownHours,
		MaxKnowledgeBaseSize:      config.MaxKnowledgeBaseSize,
		GossipMinExchangePerEvent: config.GossipMinExchangePerEvent,
		GossipMaxExchangePerEvent: config.GossipMaxExchangePerEvent,
	}, nil
}
