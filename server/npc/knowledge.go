package npc

import "time"

type KnowledgeType string

const KnowledgeAvistamentoNPC KnowledgeType = "AVISTAMENTO_NPC"

type Knowledge struct {
	Type        KnowledgeType
	EntityId    string    // id do NPC visto (no caso de AVISTAMENTO_NPC)
	BlockId     string    // onde foi visto
	PoiId       string    // POI; "" = outdoor virtual
	FirstSeenAt time.Time // primeiro avistamento na chave dedup
	LastSeenAt  time.Time // último avistamento
	SeenCount   int       // quantas vezes (avistamentos diretos)
	LearnedAt   time.Time // quando ESTE npc soube pela primeira vez
	Source      string    // "direct" ou npc_id da fonte
	Important   bool      // isenção de expiração e FIFO
}

func NewKnowledgeAvistamentoNPC(entityId, blockId, poiId, source string) Knowledge {
	now := time.Now()
	return Knowledge{
		Type:        KnowledgeAvistamentoNPC,
		EntityId:    entityId,
		BlockId:     blockId,
		PoiId:       poiId,
		FirstSeenAt: now,
		LastSeenAt:  now,
		SeenCount:   1,
		LearnedAt:   now,
		Source:      source,
		Important:   false,
	}
}
