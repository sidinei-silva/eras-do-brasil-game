package npc

import (
	"log/slog"
	"math/rand"
	"time"

	"github.com/sidinei-silva/eras-do-brasil-game/server/config"
)

type Role string

const (
	RoleMerchant   Role = "merchant"
	RoleQuestGiver Role = "quest_giver"
	RoleGuard      Role = "guard"
	RoleVillager   Role = "villager"
)

type Npc struct {
	Id                string           // Identificador único do NPC
	Name              string           // Nome do NPC
	Role              Role             // Função ou profissão do NPC (ex.: comerciante, entregador de missão, guarda)
	CurrentBlock      string           // Bloco atual onde o NPC está localizado
	CurrentActivity   Activity         // Atividade atual em que o NPC está envolvido
	Description       string           // Descrição da aparência e personalidade do NPC
	Backstory         string           // História de fundo do NPC
	Schedule          []ScheduleAction // Agenda diária das atividades do NPC
	Needs             Need             // Necessidades atuais do NPC
	HomePoi           string           // Local onde o NPC mora para dormir
	EatingPoi         string           // Local onde o NPC come
	ActivityStartedAt time.Time        // Hora em que a atividade atual começou
	NeedsWeight       NeedWeight       // Pesos para cada necessidade, usados na decisão de atividades
	CurrentPoi        string           // Ponto de interesse onde o npc está atualmente
	validPois         map[string]struct{}
	KnowledgeBase     []Knowledge          // Base de conhecimento do NPC, contendo informações que ele pode compartilhar ou usar para tomar decisões
	LastGossipedWith  map[string]time.Time // Mapa para rastrear a última vez que o NPC conversou com outros NPCs ou jogadores, para evitar repetição excessiva de fofocas
}

func NewNpc(id string, name string, role Role, currentBlock string, description string, backstory string, schedule []ScheduleAction, homePoi string, eatingPoi string, validPois []string, needsWeight NeedWeight) *Npc {
	poiSet := make(map[string]struct{}, len(validPois))
	for _, poiId := range validPois {
		poiSet[poiId] = struct{}{}
	}

	return &Npc{
		Id:              id,
		Name:            name,
		Role:            role,
		CurrentBlock:    currentBlock,
		CurrentActivity: "idle", // A atividade padrão é idle até que a agenda seja avaliada
		Description:     description,
		Backstory:       backstory,
		Schedule:        schedule,
		Needs: Need{
			Hunger:     0,
			Fatigue:    0,
			Loneliness: 0,
		},
		HomePoi:          homePoi,
		EatingPoi:        eatingPoi,
		NeedsWeight:      needsWeight,
		CurrentPoi:       "",
		validPois:        poiSet,
		KnowledgeBase:    []Knowledge{},
		LastGossipedWith: make(map[string]time.Time),
	}
}

func (npc *Npc) hasValidPoi(poiId string) bool {
	if poiId == "" {
		return true
	}

	_, exists := npc.validPois[poiId]
	return exists
}

// addOrUpdateKnowledge adiciona um novo conhecimento à base de conhecimento do NPC ou atualiza um conhecimento existente se já houver um correspondente
func (npc *Npc) addOrUpdateKnowledge(newKnowledge Knowledge, gameTime time.Time) {
	// Verifica se já existe um conhecimento correspondente na base de conhecimento
	if config.Log.NPCKnowledge {
		slog.Debug("Verificando se conhecimento existe no NPC", "npc_id", npc.Id, "knowledge", newKnowledge)
	}
	for i, k := range npc.KnowledgeBase {
		if k.Type == newKnowledge.Type && k.EntityId == newKnowledge.EntityId && k.BlockId == newKnowledge.BlockId && k.PoiId == newKnowledge.PoiId {
			// Atualiza o conhecimento existente
			npc.KnowledgeBase[i].LastSeenAt = gameTime
			npc.KnowledgeBase[i].SeenCount++
			if config.Log.NPCKnowledge {
				slog.Debug("Conhecimento atualizado no NPC", "npc_id", npc.Id, "knowledge", npc.KnowledgeBase[i])
			}
			return
		}
	}
	// Adiciona novo conhecimento
	npc.KnowledgeBase = append(npc.KnowledgeBase, newKnowledge)
	if config.Log.NPCKnowledge {
		slog.Debug("Conhecimento adicionado ao NPC", "npc_id", npc.Id, "knowledge", newKnowledge)
	}
}

// getKnowledge retorna o conhecimento correspondente ao tipo, entidade, bloco e poi especificados, se existir
func (npc *Npc) getKnowledge(knowledgeType KnowledgeType, entityId, blockId, poiId string) (Knowledge, bool) {
	for _, k := range npc.KnowledgeBase {
		if k.Type == knowledgeType && k.EntityId == entityId && k.BlockId == blockId && k.PoiId == poiId {
			return k, true
		}
	}
	return Knowledge{}, false
}

func (npc *Npc) removeExpiredKnowledge(gameTime time.Time, expirationHours map[KnowledgeType]int) {
	filtered := make([]Knowledge, 0, len(npc.KnowledgeBase))
	for _, k := range npc.KnowledgeBase {
		expiration, exists := expirationHours[k.Type]
		if !exists {
			slog.Warn("Sem config de expiração", "knowledge_type", k.Type)
			filtered = append(filtered, k) // preserva (não sei se devia, decisão)
			continue
		}
		expired := gameTime.Sub(k.LastSeenAt) > time.Duration(expiration)*time.Hour && !k.Important
		if expired {
			if config.Log.NPCKnowledge {
				slog.Debug("Conhecimento expirado removido", "npc_id", npc.Id, "knowledge", k)
			}
			continue
		}
		filtered = append(filtered, k)
	}
	npc.KnowledgeBase = filtered
}

func (npc *Npc) eligibleKnowledgeForGossip(receptorId string) []Knowledge {
	eligible := make([]Knowledge, 0, len(npc.KnowledgeBase))
	for _, k := range npc.KnowledgeBase {
		if k.EntityId != receptorId && k.Source != receptorId {
			eligible = append(eligible, k)
		}
	}
	return eligible
}

func pickRandomKnowledge(entries []Knowledge, minCount, maxCount int) []Knowledge {
	if len(entries) == 0 {
		return nil
	}
	if minCount < 1 {
		minCount = 1
	}
	if maxCount < minCount {
		maxCount = minCount
	}

	count := minCount
	if maxCount > minCount {
		count = minCount + rand.Intn(maxCount-minCount+1)
	}
	if count > len(entries) {
		count = len(entries)
	}

	perm := rand.Perm(len(entries))
	selected := make([]Knowledge, 0, count)
	for i := 0; i < count; i++ {
		selected = append(selected, entries[perm[i]])
	}
	return selected
}

func mergeGossipedKnowledge(receptor *Npc, source *Npc, knowledge Knowledge, gameTime time.Time) {
	if existingK, found := receptor.getKnowledge(knowledge.Type, knowledge.EntityId, knowledge.BlockId, knowledge.PoiId); found {
		if knowledge.LastSeenAt.After(existingK.LastSeenAt) {
			for i, k := range receptor.KnowledgeBase {
				if k.Type == knowledge.Type && k.EntityId == knowledge.EntityId && k.BlockId == knowledge.BlockId && k.PoiId == knowledge.PoiId {
					receptor.KnowledgeBase[i].LastSeenAt = knowledge.LastSeenAt
					if config.Log.NPCKnowledge {
						slog.Debug("Conhecimento atualizado via fofoca",
							"receptor_id", receptor.Id,
							"emissor_id", source.Id,
							"knowledge_entity_id", knowledge.EntityId,
							"last_seen_at", knowledge.LastSeenAt,
						)
					}
					break
				}
			}
		}
		return
	}

	newK := Knowledge{
		Type:        knowledge.Type,
		EntityId:    knowledge.EntityId,
		BlockId:     knowledge.BlockId,
		PoiId:       knowledge.PoiId,
		FirstSeenAt: knowledge.FirstSeenAt,
		LastSeenAt:  knowledge.LastSeenAt,
		SeenCount:   knowledge.SeenCount,
		LearnedAt:   gameTime,
		Source:      source.Id,
		Important:   knowledge.Important,
	}
	receptor.KnowledgeBase = append(receptor.KnowledgeBase, newK)
	if config.Log.NPCKnowledge {
		slog.Debug("Conhecimento adicionado via fofoca",
			"receptor_id", receptor.Id,
			"emissor_id", source.Id,
			"knowledge_entity_id", newK.EntityId,
			"learned_at", newK.LearnedAt,
		)
	}
}

// exchangeKnowledgeWith executa a troca bidirecional em uma única passada.
// Isso evita que o par seja reprocessado com KB já mutada pelo primeiro lado.
func (npc *Npc) exchangeKnowledgeWith(otherNpc *Npc, gameTime time.Time, knowledgeConfig KnowledgeConfig) bool {
	firstSide := pickRandomKnowledge(npc.eligibleKnowledgeForGossip(otherNpc.Id), knowledgeConfig.GossipMinExchangePerEvent, knowledgeConfig.GossipMaxExchangePerEvent)
	secondSide := pickRandomKnowledge(otherNpc.eligibleKnowledgeForGossip(npc.Id), knowledgeConfig.GossipMinExchangePerEvent, knowledgeConfig.GossipMaxExchangePerEvent)

	if len(firstSide) == 0 && len(secondSide) == 0 {
		return false
	}

	for _, knowledge := range firstSide {
		mergeGossipedKnowledge(otherNpc, npc, knowledge, gameTime)
	}
	for _, knowledge := range secondSide {
		mergeGossipedKnowledge(npc, otherNpc, knowledge, gameTime)
	}

	npc.LastGossipedWith[otherNpc.Id] = gameTime
	otherNpc.LastGossipedWith[npc.Id] = gameTime

	return true
}
